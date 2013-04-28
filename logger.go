package tigertonic

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Logger struct {
	*log.Logger
	handler http.Handler
}

func Logged(handler http.Handler) *Logger {
	return &Logger{
		Logger: log.New(
			os.Stderr,
			fmt.Sprintf("%s ", requestID()),
			log.Ltime | log.Lmicroseconds,
		),
		handler: handler,
	}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.Printf("> %s %s %s\n", r.Method, r.URL, r.Proto)
	for name, values := range r.Header {
		for _, value := range values {
			l.Printf("> %s: %s\n", name, value)
		}
	}
	l.Println(">")
	r.Body = &readCloser{
		ReadCloser: r.Body,
		Logger: l.Logger,
	}
	l.handler.ServeHTTP(&responseWriter{
		ResponseWriter: w,
		Logger: l.Logger,
		request: r,
	}, r)
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
		encodingBase62[i] = s[uint8(float32(61) * float32(buf[i]) / float32(255))]
	}
}

type readCloser struct {
	io.ReadCloser
	*log.Logger
}

func (r *readCloser) Read(buf []byte) (int, error) {
	n, err := r.ReadCloser.Read(buf)
	if 0 < n && nil == err {
		r.Println(">", string(buf))
	}
	return n, err
}

func requestID() string {
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
	return string(buf)
}

type responseWriter struct {
	http.ResponseWriter
	*log.Logger
	request *http.Request
}

func (w *responseWriter) Write(buf []byte) (int, error) {
	w.Println("<", string(buf))
	return w.ResponseWriter.Write(buf)
}

func (w *responseWriter) WriteHeader(status int) {
	w.Printf("< %s %d %s\n", w.request.Proto, status, http.StatusText(status))
	for name, values := range w.Header() {
		for _, value := range values {
			w.Printf("< %s: %s\n", name, value)
		}
	}
	w.Println("<")
}
