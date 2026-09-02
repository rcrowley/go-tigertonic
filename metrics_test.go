package tigertonic

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"go.withmatt.com/metrics"
)

func TestCounter(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	counter := Counted(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), "counted", metrics.NewSet())
	counter.ServeHTTP(w, r)
	if 1 != counter.Get() {
		t.Fatal(counter.Get())
	}
}

func TestCounterByStatus(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	counterByStatus := CountedByStatus(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), "counted_by_status", metrics.NewSet())
	counterByStatus.ServeHTTP(w, r)
	if 1 != counterByStatus.counters[200].Get() {
		t.Fatal(counterByStatus.counters[200].Get())
	}
	if 0 != counterByStatus.counters[500].Get() {
		t.Fatal(counterByStatus.counters[500].Get())
	}
}

func TestCounterByStatusXX(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	counterByStatusXX := CountedByStatusXX(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), "counted_by_status_xx", metrics.NewSet())
	counterByStatusXX.ServeHTTP(w, r)
	if 1 != counterByStatusXX.counter2xx.Get() {
		t.Fatal(counterByStatusXX.counter2xx.Get())
	}
	if 0 != counterByStatusXX.counter5xx.Get() {
		t.Fatal(counterByStatusXX.counter5xx.Get())
	}
}

func TestTimer(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	set := metrics.NewSet()
	timer := Timed(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), "timed", set)
	timer.ServeHTTP(w, r)
	var buf bytes.Buffer
	if _, err := set.WritePrometheus(&buf); nil != err {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "timed_count 1") {
		t.Fatal(buf.String())
	}
}
