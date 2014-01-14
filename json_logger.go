package tigertonic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// JSONLogged is an http.Handler that logs requests and responses in a parse-able json line.
// complete with paths, statuses, headers, and bodies.  Sensitive information may be
// redacted by a user-defined function.
type JSONLogger struct {
	logger           *log.Logger
	handler          http.Handler
	redactor         Redactor
	RequestIDCreator RequestIDCreator
}

// JSONLogged returns an http.Handler that logs requests and responses in a parse-able json line.
// complete with paths, statuses, headers, and bodies.  Sensitive information may be
// redacted by a user-defined function.
func JSONLogged(handler http.Handler, redactor Redactor) *JSONLogger {
	return &JSONLogger{
		logger:           log.New(os.Stdout, "", 0),
		handler:          handler,
		redactor:         redactor,
		RequestIDCreator: requestIDCreator,
	}
}

// ServeHTTP wraps the http.Request and http.ResponseWriter to capture the input and output for logging
func (jl *JSONLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := jl.RequestIDCreator(r)
	startTime := time.Now()

	teeWriter := NewTeeResponseWriter(w)
	body := &jsonReadCloser{r.Body, bytes.Buffer{}}
	r.Body = body

	jl.handler.ServeHTTP(teeWriter, r)

	message := fmt.Sprintf(
		"%s %s %s\n%s %d %s",
		r.Method,
		r.URL.RequestURI(),
		r.Proto,
		r.Proto,
		teeWriter.StatusCode,
		http.StatusText(teeWriter.StatusCode),
	)

	outObject := &logObject{
		message,
		"http",
		requestID,
		int64(time.Since(startTime) / time.Millisecond),
		httpObject{
			fmt.Sprintf("%d.%d", r.ProtoMajor, r.ProtoMinor),
			buildRequest(r, body),
			buildResponse(teeWriter),
		},
	}

	out, err := json.Marshal(outObject)
	if err != nil {
		jl.logger.Println("Error while formatting request log:", err)
		return
	}

	outString := string(out)
	if nil != jl.redactor {
		outString = jl.redactor(outString)
	}

	jl.logger.Println(outString)
}

// Builds up the request object that will appear inside http.request in the logs
func buildRequest(r *http.Request, b *jsonReadCloser) httpRequestObject {
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, "; ")
	}

	return httpRequestObject{
		r.Method,
		r.URL.RequestURI(),
		headers,
		b.Bytes.String(),
	}
}

// Builds up the response object that will appear inside http.response in the logs
func buildResponse(tw *TeeResponseWriter) httpResponseObject {
	headers := make(map[string]string)
	for name, values := range tw.Header() {
		headers[name] = strings.Join(values, "; ")
	}

	return httpResponseObject{
		tw.StatusCode,
		http.StatusText(tw.StatusCode),
		headers,
		string(tw.Body.String()),
	}
}

type logObject struct {
	Message   string     `json:"@message"`
	Type      string     `json:"@type"`
	RequestId RequestID  `json:"@request_id"`
	Duration  int64      `json:"duration"`
	Http      httpObject `json:"http"`
}

type httpObject struct {
	Version  string             `json:"version"`
	Request  httpRequestObject  `json:"request"`
	Response httpResponseObject `json:"response"`
}

type httpRequestObject struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type httpResponseObject struct {
	Status  int               `json:"status"`
	Reason  string            `json:"reason"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type jsonReadCloser struct {
	io.ReadCloser
	Bytes bytes.Buffer
}

func (r *jsonReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	r.Bytes.Write(p[:n])

	return n, err
}
