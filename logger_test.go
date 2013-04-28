package tigertonic

import (
	"bytes"
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
	if `> POST /foo HTTP/1.1
> Accept: application/json
> Content-Type: application/json
>
> {"foo":"bar"}
< HTTP/1.1 200 OK
<
< {"foo":"bar"}
` != b.String() {
		t.Fatal(b.String())
	}
}
