package tigertonic

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestLogger(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	logger := Logged(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}))
	b := &bytes.Buffer{}
	logger.Logger = log.New(b, "", 0)
	logger.ServeHTTP(w, r)
	requestID := b.String()[:16]
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
	) != b.String() {
		t.Fatal(b.String())
	}
}
