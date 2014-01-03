package tigertonic

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"log"
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
func Counted(
	handler http.Handler,
	name string,
	registry metrics.Registry,
) *Counter {
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

// StatusCodeRememberingWriter is an impl of http.ResponseWriter which keeps
// track of the status code for further processing
type StatusCodeRememberingWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader simply wraps the underlying WriteHeader method, storing status
func (self *StatusCodeRememberingWriter) WriteHeader(status int) {
	self.StatusCode = status
	self.ResponseWriter.WriteHeader(status)
}

// StatusCodeCounter uses wrap.Processed to keep track of http status codes
type StatusCodeCounter struct {
	counter1xx metrics.Counter
	counter2xx metrics.Counter
	counter3xx metrics.Counter
	counter4xx metrics.Counter
	counter5xx metrics.Counter
	handler    http.Handler
}

func (self *StatusCodeCounter) WrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &StatusCodeRememberingWriter{ResponseWriter: w}
}

func (self *StatusCodeCounter) Process(w http.ResponseWriter, r *http.Request) {
	var wrapped *StatusCodeRememberingWriter
	var ok bool
	if wrapped, ok = w.(*StatusCodeRememberingWriter); !ok {
		log.Printf("StatusCodeCounter was not able to properly cast ResponseWriter to extract response code")
		return
	}

	if wrapped.StatusCode < 200 {
		self.counter1xx.Inc(1)
	} else if wrapped.StatusCode < 300 {
		self.counter2xx.Inc(1)
	} else if wrapped.StatusCode < 400 {
		self.counter3xx.Inc(1)
	} else if wrapped.StatusCode < 500 {
		self.counter4xx.Inc(1)
	} else {
		self.counter5xx.Inc(1)
	}
}

func StatusCodeCounted(handler http.Handler, name string, registry metrics.Registry) http.Handler {
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	counters := make(map[string]metrics.Counter)
	counters["1xx"] = metrics.NewCounter()
	counters["2xx"] = metrics.NewCounter()
	counters["3xx"] = metrics.NewCounter()
	counters["4xx"] = metrics.NewCounter()
	counters["5xx"] = metrics.NewCounter()
	for code, counter := range counters {
		registry.Register(fmt.Sprintf("%s-%s", name, code), counter)
	}

	return Processed(handler, &StatusCodeCounter{
		counter1xx: counters["1xx"],
		counter2xx: counters["2xx"],
		counter3xx: counters["3xx"],
		counter4xx: counters["4xx"],
		counter5xx: counters["5xx"],
		handler:    handler,
	})
}
