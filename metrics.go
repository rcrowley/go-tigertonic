package tigertonic

import (
	"github.com/rcrowley/go-metrics"
	"net/http"
	"time"
)

type Counter struct {
	metrics.Counter
	handler http.Handler
}

func Counted(handler http.Handler) *Counter {
	return &Counter{
		Counter: metrics.NewCounter(),
		handler: handler,
	}
	// TODO Register
}

func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
	c.Inc(1)
}

type Timer struct {
	metrics.Timer
	handler http.Handler
}

func Timed(handler http.Handler) *Timer {
	return &Timer{
		Timer:   metrics.NewTimer(),
		handler: handler,
	}
	// TODO Register
}

func (t *Timer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer t.UpdateSince(time.Now())
	t.handler.ServeHTTP(w, r)
}
