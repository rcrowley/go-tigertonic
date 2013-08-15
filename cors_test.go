package tigertonic

import (
	"net/http"
	"net/url"
	"testing"
)

type TestResponse struct {
	ImportantInfo string `json:"important_info"`
}

// GET /baz
func get(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *TestResponse, error) {
	return http.StatusOK, nil, &TestResponse{"i love you"}, nil
}

func TestCORSOrigin(t *testing.T) {
	mux := NewTrieServeMux()
	mux.Handle("GET", "/foo", NewCORSBuilder().SetAllowedOrigin("*").Build(Marshaled(get)))
	mux.Handle("GET", "/baz", NewCORSBuilder().SetAllowedOrigin("http://gooddomain.com").Build(Marshaled(get)))

	w := &testResponseWriter{}
	r, _ := http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set("Access-Control-Request-Method", "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "GET, HEAD, OPTIONS" != w.Header().Get("Access-Control-Allow-Methods") {
		t.Fatal(w.Header().Get("Allow"))
	}

	// requesting secured resource with invalid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set("Origin", "http://baddomain.com")
	r.Header.Set("Access-Control-Request-Method", "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "null" != w.Header().Get("Access-Control-Allow-Origin") {
		t.Fatal(w.Header().Get("Access-Control-Allow-Origin"))
	}

	// requesting unsecured/wildcard resource with invalid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	r.Header.Set("Origin", "http://baddomain.com")
	r.Header.Set("Access-Control-Request-Method", "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "*" != w.Header().Get("Access-Control-Allow-Origin") {
		t.Fatal(w.Header().Get("Access-Control-Allow-Origin"))
	}

	// requesting secured resource with valid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set("Origin", "http://gooddomain.com")
	r.Header.Set("Access-Control-Request-Method", "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "http://gooddomain.com" != w.Header().Get("Access-Control-Allow-Origin") {
		t.Fatal(w.Header().Get("Access-Control-Allow-Origin"))
	}
}
