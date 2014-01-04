package tigertonic

import (
	"bytes"
	"net/http"
)

// TeeResponseWriter is an http.ResponseWriter that both writes and records the
// response for post-processing.
type TeeResponseWriter struct {
	http.ResponseWriter
	Body   bytes.Buffer
	Status int
}

// NewTeeResponseWriter constructs a new TeeResponseWriter that wraps another
// http.ResponseWriter.
func NewTeeResponseWriter(w http.ResponseWriter) *TeeResponseWriter {
	return &TeeResponseWriter{ResponseWriter: w}
}

// Write writes the byte slice to the client via the underlying
// http.ResponseWriter and records it for post-processing.
func (w *TeeResponseWriter) Write(p []byte) (int, error) {
	if n, err := w.ResponseWriter.Write(p); nil != err {
		return n, err
	}
	return w.Body.Write(p)
}

// WriteHeader writes the response line and headers to the client via the
// underlying http.ResponseWriter and records the status for post-processing.
func (w *TeeResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.Status = status
}
