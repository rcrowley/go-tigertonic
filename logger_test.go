package tigertonic

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	logger := Logged(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), nil)
	b := &bytes.Buffer{}
	logger.Logger = log.New(b, "", 0)
	logger.ServeHTTP(w, r)
	s := b.String()
	requestID := s[:16]
	if fmt.Sprintf(
		`%s > POST /foo HTTP/1.1
%s > Accept: application/json
%s > Content-Type: application/json
%s >
%s > {"foo":"bar"}
%s < HTTP/1.1 200 OK
%s < Content-Type: application/json
%s <
%s < {"foo":"bar"}
`,
		requestID,
		requestID,
		requestID,
		requestID,
		requestID,
		requestID,
		requestID,
		requestID,
		requestID,
	) != s {
		t.Fatal(s)
	}
}

func TestRedactor(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	logger := Logged(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"SECRET"}, nil
	}), func(s string) string {
		return strings.Replace(s, "SECRET", "REDACTED", -1)
	})
	b := &bytes.Buffer{}
	logger.Logger = log.New(b, "", 0)
	logger.ServeHTTP(w, r)
	s := b.String()
	if strings.Contains(s, "SECRET") {
		t.Fatal(s)
	}
	if !strings.Contains(s, "REDACTED") {
		t.Fatal(s)
	}
}
