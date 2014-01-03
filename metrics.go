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

// CounterByStatus is an http.Handler that counts responses by their HTTP
// status via go-metrics.
type CounterByStatus struct {
	counters map[int]metrics.Counter
	handler  http.Handler
}

// CountedByStatus returns an http.Handler that passes requests to an
// underlying http.Handler and then counts the response by its HTTP status via
// go-metrics.
func CountedByStatus(
	handler http.Handler,
	name string,
	registry metrics.Registry,
) *CounterByStatus {
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	counters := map[int]metrics.Counter{
		100: metrics.NewCounter(),
		101: metrics.NewCounter(),
		200: metrics.NewCounter(),
		201: metrics.NewCounter(),
		202: metrics.NewCounter(),
		203: metrics.NewCounter(),
		204: metrics.NewCounter(),
		205: metrics.NewCounter(),
		206: metrics.NewCounter(),
		300: metrics.NewCounter(),
		301: metrics.NewCounter(),
		302: metrics.NewCounter(),
		303: metrics.NewCounter(),
		304: metrics.NewCounter(),
		305: metrics.NewCounter(),
		306: metrics.NewCounter(),
		307: metrics.NewCounter(),
		400: metrics.NewCounter(),
		401: metrics.NewCounter(),
		402: metrics.NewCounter(),
		403: metrics.NewCounter(),
		404: metrics.NewCounter(),
		405: metrics.NewCounter(),
		406: metrics.NewCounter(),
		407: metrics.NewCounter(),
		408: metrics.NewCounter(),
		409: metrics.NewCounter(),
		410: metrics.NewCounter(),
		411: metrics.NewCounter(),
		412: metrics.NewCounter(),
		413: metrics.NewCounter(),
		414: metrics.NewCounter(),
		415: metrics.NewCounter(),
		416: metrics.NewCounter(),
		417: metrics.NewCounter(),
		500: metrics.NewCounter(),
		501: metrics.NewCounter(),
		502: metrics.NewCounter(),
		503: metrics.NewCounter(),
		504: metrics.NewCounter(),
		505: metrics.NewCounter(),
	}
	for status, counter := range counters {
		registry.Register(fmt.Sprintf("%s-%d", name, status), counter)
	}
	return &CounterByStatus{
		counters: counters,
		handler:  handler,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the response by its HTTP status via go-metrics.
func (c *CounterByStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srw := &statusResponseWriter{ResponseWriter: w}
	c.handler.ServeHTTP(srw, r)
	c.counters[srw.Status].Inc(1)
}

// CounterByStatusXX is an http.Handler that counts responses by the first
// digit of their HTTP status via go-metrics.
type CounterByStatusXX struct {
	counter1xx, counter2xx, counter3xx, counter4xx, counter5xx metrics.Counter
	handler  http.Handler
}

// CountedByStatusXX returns an http.Handler that passes requests to an
// underlying http.Handler and then counts the response by the first digit of
// its HTTP status via go-metrics.
func CountedByStatusXX(
	handler http.Handler,
	name string,
	registry metrics.Registry,
) *CounterByStatusXX {
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	c := &CounterByStatusXX{
		counter1xx: metrics.NewCounter(),
		counter2xx: metrics.NewCounter(),
		counter3xx: metrics.NewCounter(),
		counter4xx: metrics.NewCounter(),
		counter5xx: metrics.NewCounter(),
		handler:    handler,
	}
	registry.Register(fmt.Sprintf("%s-1xx", name), c.counter1xx)
	registry.Register(fmt.Sprintf("%s-2xx", name), c.counter2xx)
	registry.Register(fmt.Sprintf("%s-3xx", name), c.counter3xx)
	registry.Register(fmt.Sprintf("%s-4xx", name), c.counter4xx)
	registry.Register(fmt.Sprintf("%s-5xx", name), c.counter5xx)
	return c
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the response by its HTTP status via go-metrics.
func (c *CounterByStatusXX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srw := &statusResponseWriter{ResponseWriter: w}
	c.handler.ServeHTTP(srw, r)
	if srw.Status < 200 {
		c.counter1xx.Inc(1)
	} else if srw.Status < 300 {
		c.counter2xx.Inc(1)
	} else if srw.Status < 400 {
		c.counter3xx.Inc(1)
	} else if srw.Status < 500 {
		c.counter4xx.Inc(1)
	} else {
		c.counter5xx.Inc(1)
	}
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

type statusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.Status = status
}
