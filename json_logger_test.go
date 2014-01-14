package tigertonic

import (
	"testing"
	"net/url"
	"net/http"
	"bytes"
	"log"
	"encoding/json"
	"reflect"
	"strings"
)

func TestJSONLogger(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest(
		"POST",
		"http://example.com/foo?bar=baz",
		bytes.NewBufferString(`{"foo":"bar"}`),
	)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")

	logger := JSONLogged(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), nil)

	logger.RequestIDCreator = func (r *http.Request) RequestID {
		return "request-id"
	}

	b := &bytes.Buffer{}
	logger.logger = log.New(b, "", 0)
	logger.ServeHTTP(w, r)

	var m logObject

	err := json.Unmarshal(b.Bytes(), &m)
	if (err != nil) {
		t.Fatal(err)
	}

	requestHeaders := make(map[string]string)
	requestHeaders["accept"] = "application/json"
	requestHeaders["content-type"] = "application/json"

	responseHeaders := make(map[string]string)
	responseHeaders["content-type"] = "application/json"

	expected := logObject{
		"POST /foo?bar=baz HTTP/1.1\nHTTP/1.1 200 OK",
		"http",
		"request-id",
		int64(0),
		httpObject{
			"1.1",
			httpRequestObject{
				"POST",
				"/foo?bar=baz",
				requestHeaders,
				"{\"foo\":\"bar\"}",
			},
			httpResponseObject{
				200,
				"OK",
				responseHeaders,
				"{\"foo\":\"bar\"}\n",
			},
		},
	}

	if reflect.DeepEqual(expected, m) == false {
		t.Fatalf("Log object was incorrect\nExpected\n%+v\nGot\n%+v", expected, m)
	}
}

func TestJSONLoggerRedactor(t *testing.T) {
	w := &testResponseWriter{}

	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")

	logger := JSONLogged(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
			return http.StatusOK, nil, &testResponse{"SECRET"}, nil
		}), func(s string) string {
			return strings.Replace(s, "SECRET", "REDACTED", -1)
		})

	b := &bytes.Buffer{}
	logger.logger = log.New(b, "", 0)
	logger.ServeHTTP(w, r)
	s := b.String()

	if strings.Contains(s, "SECRET") {
		t.Fatal(s)
	}

	if !strings.Contains(s, "REDACTED") {
		t.Fatal(s)
	}
}
