package tigertonic

import "net/http"

type first []http.Handler

// First returns an http.Handler that, for each handler in its slice of
// handlers, calls ServeHTTP until the first one that calls w.WriteHeader.
func First(handlers ...http.Handler) first {
	return handlers
}

func (f first) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w0 := &firstResponseWriter{w, false}
	for _, h := range f {
		h.ServeHTTP(w0, r)
		if w0.written {
			break
		}
	}
}

type firstResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *firstResponseWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}
