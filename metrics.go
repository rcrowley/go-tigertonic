package tigertonic

import (
	"net/http"
)

type Counter struct {
	h http.Handler
}

func Counted(h http.Handler) *Counter {
	return &Counter{h}
}

func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.h.ServeHTTP(w, r)
	// TODO go-metrics
}

type Timer struct {
	h http.Handler
}

func Timed(h http.Handler) *Timer {
	return &Timer{h}
}

func (t *Timer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO go-metrics
	t.h.ServeHTTP(w, r)
	// TODO go-metrics
}
