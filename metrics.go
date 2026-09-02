package tigertonic

import (
	"fmt"
	"net/http"
	"time"

	"go.withmatt.com/metrics"
)

var counterStatusCodes = []int{
	100, 101,
	200, 201, 202, 203, 204, 205, 206,
	300, 301, 302, 303, 304, 305, 306, 307,
	400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 422,
	500, 501, 502, 503, 504, 505,
}

// Counter is an http.Handler that counts requests via go-metrics.
type Counter struct {
	*metrics.Uint64
	handler http.Handler
}

// Counted returns an http.Handler that passes requests to an underlying
// http.Handler and then counts the request via go-metrics.
func Counted(handler http.Handler, name string, set *metrics.Set) *Counter {
	return &Counter{
		Uint64:  newCounter(set, name),
		handler: handler,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the request via go-metrics.
func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
	c.Inc()
}

// CounterByStatus is an http.Handler that counts responses by their HTTP
// status code via go-metrics.
type CounterByStatus struct {
	counters map[int]*metrics.Uint64
	handler  http.Handler
}

// CountedByStatus returns an http.Handler that passes requests to an
// underlying http.Handler and then counts the response by its HTTP status code
// via go-metrics.
func CountedByStatus(
	handler http.Handler,
	name string,
	set *metrics.Set,
) *CounterByStatus {
	counters := make(map[int]*metrics.Uint64, len(counterStatusCodes))
	for _, code := range counterStatusCodes {
		counters[code] = newCounter(
			set,
			name,
			"status", fmt.Sprintf("%d", code),
		)
	}
	return &CounterByStatus{
		counters: counters,
		handler:  handler,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the response by its HTTP status code via go-metrics.
func (c *CounterByStatus) ServeHTTP(w0 http.ResponseWriter, r *http.Request) {
	w := NewTeeHeaderResponseWriter(w0)
	c.handler.ServeHTTP(w, r)
	c.counters[w.StatusCode].Inc()
}

// CounterByStatusXX is an http.Handler that counts responses by the first
// digit of their HTTP status code via go-metrics.
type CounterByStatusXX struct {
	counter1xx, counter2xx, counter3xx, counter4xx, counter5xx *metrics.Uint64
	handler                                                    http.Handler
}

// CountedByStatusXX returns an http.Handler that passes requests to an
// underlying http.Handler and then counts the response by the first digit of
// its HTTP status code via go-metrics.
func CountedByStatusXX(
	handler http.Handler,
	name string,
	set *metrics.Set,
) *CounterByStatusXX {
	return &CounterByStatusXX{
		counter1xx: newCounter(set, name, "status", "1xx"),
		counter2xx: newCounter(set, name, "status", "2xx"),
		counter3xx: newCounter(set, name, "status", "3xx"),
		counter4xx: newCounter(set, name, "status", "4xx"),
		counter5xx: newCounter(set, name, "status", "5xx"),
		handler:    handler,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the response by its HTTP status code via go-metrics.
func (c *CounterByStatusXX) ServeHTTP(w0 http.ResponseWriter, r *http.Request) {
	w := NewTeeHeaderResponseWriter(w0)
	c.handler.ServeHTTP(w, r)
	if w.StatusCode < 200 {
		c.counter1xx.Inc()
	} else if w.StatusCode < 300 {
		c.counter2xx.Inc()
	} else if w.StatusCode < 400 {
		c.counter3xx.Inc()
	} else if w.StatusCode < 500 {
		c.counter4xx.Inc()
	} else {
		c.counter5xx.Inc()
	}
}

// Timer is an http.Handler that times requests via go-metrics.
type Timer struct {
	*metrics.Histogram
	handler http.Handler
}

// Timed returns an http.Handler that passes requests to an underlying
// http.Handler and then updates a histogram with the duration of the request
// via go-metrics.
func Timed(handler http.Handler, name string, set *metrics.Set) *Timer {
	return &Timer{
		Histogram: newHistogram(set, name),
		handler:   handler,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then
// updates the histogram with the duration of the request via go-metrics.
func (t *Timer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer t.UpdateDuration(time.Now())
	t.handler.ServeHTTP(w, r)
}

func newCounter(set *metrics.Set, family string, tags ...string) *metrics.Uint64 {
	if nil == set {
		return metrics.NewCounter(family, tags...)
	}
	return set.NewCounter(family, tags...)
}

func newHistogram(set *metrics.Set, family string, tags ...string) *metrics.Histogram {
	if nil == set {
		return metrics.NewHistogram(family, tags...)
	}
	return set.NewHistogram(family, tags...)
}
