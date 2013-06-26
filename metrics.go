package tigertonic

import (
	"github.com/rcrowley/go-metrics"
	"net/http"
	"time"
)

// Counter is an http.Handler that counts requests via go-metrics.
type Counter struct {
	metrics.Counter
	handler http.Handler
}

// Counted returns an http.Handler that passes requests to an underlying
// http.Handler and then counts the request via go-metrics.
func Counted(handler http.Handler, name string, registry metrics.Registry) *Counter {
	counter := &Counter{
		Counter: metrics.NewCounter(),
		handler: handler,
	}
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	registry.Register(name, counter)
	return counter
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the request via go-metrics.
func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
	c.Inc(1)
}

// Timer is an http.Handler that counts requests via go-metrics.
type Timer struct {
	metrics.Timer
	handler http.Handler
}

// Timed returns an http.Handler that starts a timer, passes requests to an
// underlying http.Handler, stops the timer, and updates the timer via
// go-metrics.
func Timed(handler http.Handler, name string, registry metrics.Registry) *Timer {
	timer := &Timer{
		Timer:   metrics.NewTimer(),
		handler: handler,
	}
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	registry.Register(name, timer)
	return timer
}

// ServeHTTP starts a timer, passes the request to the underlying http.Handler,
// stops the timer, and updates the timer via go-metrics.
func (t *Timer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer t.UpdateSince(time.Now())
	t.handler.ServeHTTP(w, r)
}
