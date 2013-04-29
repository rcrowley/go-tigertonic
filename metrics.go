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

func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
	c.Inc(1)
}

type Timer struct {
	metrics.Timer
	handler http.Handler
}

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

func (t *Timer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer t.UpdateSince(time.Now())
	t.handler.ServeHTTP(w, r)
}
