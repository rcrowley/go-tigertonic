// Package mocking makes testing Tiger Tonic services easier.
package mocking

import (
	"net/http"
	"net/url"
)

func Header(h http.Header) http.Header {
	h0 := make(http.Header)
	h0.Add("Accept", "application/json")
	h0.Add("Content-Type", "application/json")
	if nil != h {
		for key, values := range h {
			for _, value := range values {
				h0.Add(key, value)
			}
		}
	}
	return h0
}

type TestableHandler interface {
	Handler(*http.Request) (http.Handler, string)
}

func URL(h TestableHandler, method, rawurl string) *url.URL {
	u, err := url.ParseRequestURI(rawurl)
	if nil != err {
		panic(err)
	}
	if nil != h {
		rq := &http.Request{
			Method: method,
			URL:    u,
		}
		var ok bool
		for {
			h1, _ := h.Handler(rq)
			if h, ok = h1.(TestableHandler); !ok {
				break
			}
		}
	}
	return u
}
