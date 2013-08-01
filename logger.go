package tigertonic

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Logger is an http.Handler that logs requests and responses, complete with
// paths, statuses, headers, and bodies.  Sensitive information may be redacted
// by a user-defined function.
type Logger struct {
	*log.Logger
	handler  http.Handler
	redactor Redactor
}

// Logged returns an http.Handler that logs requests and responses, complete
// with paths, statuses, headers, and bodies.  Sensitive information may be
// redacted by a user-defined function.
func Logged(handler http.Handler, redactor Redactor) *Logger {
	return &Logger{
		Logger:   log.New(os.Stdout, "", log.Ltime|log.Lmicroseconds),
		handler:  handler,
		redactor: redactor,
	}
}

// Output overrides log.Logger's Output method, calling our redactor first.
func (l *Logger) Output(calldepth int, s string) error {
	if nil != l.redactor {
		s = l.redactor(s)
	}
	return l.Logger.Output(calldepth, s)
}

// Print is identical to log.Logger's Print but uses our overridden Output.
func (l *Logger) Print(v ...interface{}) { l.Output(2, fmt.Sprint(v...)) }

// Printf is identical to log.Logger's Print but uses our overridden Output.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

// Println is identical to log.Logger's Print but uses our overridden Output.
func (l *Logger) Println(v ...interface{}) { l.Output(2, fmt.Sprintln(v...)) }

// ServeHTTP wraps the http.Request and http.ResponseWriter to log to standard
// output and pass through to the underlying http.Handler.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := NewRequestID()
	l.Printf("%s > %s %s %s\n", requestID, r.Method, r.URL.Path, r.Proto)
	for name, values := range r.Header {
		for _, value := range values {
			l.Printf("%s > %s: %s\n", requestID, name, value)
		}
	}
	l.Println(requestID, ">")
	r.Body = &readCloser{
		ReadCloser: r.Body,
		Logger:     l,
		requestID:  requestID,
	}
	l.handler.ServeHTTP(&responseWriter{
		ResponseWriter: w,
		Logger:         l,
		request:        r,
		requestID:      requestID,
	}, r)
}

// A Redactor is a function that takes and returns a string.  It is called
// to allow sensitive information to be redacted before it is logged.
type Redactor func(string) string

// A unique RequestID is given to each request and is included with each line
// of each log entry.
type RequestID string

// NewRequestID returns a new 16-character random RequestID.
func NewRequestID() RequestID {
	buf := make([]byte, 16)
	i := 0
	for i < 16 {
		n, err := rand.Read(buf[i:])
		if nil != err {
			panic(err)
		}
		i += n
	}
	for i = 0; i < 16; i++ {
		buf[i] = encodingBase62[buf[i]]
	}
	return RequestID(buf)
}

var encodingBase62 [256]byte

func init() {
	buf := make([]byte, 256)
	i := 0
	for i < 256 {
		n, err := rand.Read(buf[i:])
		if nil != err {
			panic(err)
		}
		i += n
	}
	s := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 256; i++ {
		encodingBase62[i] = s[uint8(float32(61)*float32(buf[i])/float32(255))]
	}
}

type readCloser struct {
	io.ReadCloser
	*Logger
	requestID RequestID
}

func (r *readCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if 0 < n && nil == err {
		r.Println(r.requestID, ">", string(p[:n]))
	}
	return n, err
}

type responseWriter struct {
	http.ResponseWriter
	*Logger
	request   *http.Request
	requestID RequestID
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if '\n' == p[len(p)-1] {
		w.Println(w.requestID, "<", string(p[:len(p)-1]))
	} else {
		w.Println(w.requestID, "<", string(p))
	}
	return w.ResponseWriter.Write(p)
}

func (w *responseWriter) WriteHeader(status int) {
	w.Printf("%s < %s %d %s\n", w.requestID, w.request.Proto, status, http.StatusText(status))
	for name, values := range w.Header() {
		for _, value := range values {
			w.Printf("%s < %s: %s\n", w.requestID, name, value)
		}
	}
	w.Println(w.requestID, "<")
	w.ResponseWriter.WriteHeader(status)
}
