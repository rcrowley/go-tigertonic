package tigertonic

import (
	"log"
	"net/http"
)

type Logger struct {
	h http.Handler
}

func Logged(h http.Handler) *Logger {
	return &Logger{h}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.h.ServeHTTP(w, r) // TODO Give it a wrapped w to capture everything.
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	// TODO Full request/response logging.
}
