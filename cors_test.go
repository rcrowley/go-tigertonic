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

func TestCORSOPTIONS(t *testing.T) {
	mux := NewTrieServeMux()
	mux.Handle("GET", "/foo", NewCORSBuilder().SetAllowedOrigin("*").Build(Marshaled(get)))
	mux.Handle("GET", "/baz", NewCORSBuilder().SetAllowedOrigin("http://gooddomain.com").Build(Marshaled(get)))

	w := &testResponseWriter{}
	r, _ := http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "GET, HEAD, OPTIONS" != w.Header().Get(CORSAllowMethods) {
		t.Fatal(w.Header().Get("Allow"))
	}

	// requesting secured resource with invalid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set(CORSRequestOrigin, "http://baddomain.com")
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "null" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// requesting unsecured/wildcard resource with invalid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	r.Header.Set(CORSRequestOrigin, "http://baddomain.com")
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "*" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// requesting secured resource with valid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set(CORSRequestOrigin, "http://gooddomain.com")
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "http://gooddomain.com" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}
}

func TestCORSOrigin(t *testing.T) {
	mux := NewTrieServeMux()
	mux.Handle("GET", "/foo", NewCORSBuilder().SetAllowedOrigin("*").Build(Marshaled(get)))
	mux.Handle("GET", "/baz", NewCORSBuilder().SetAllowedOrigin("http://gooddomain.com").Build(Marshaled(get)))

	// wildcard
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set(CORSRequestOrigin, "http://gooddomain.com")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "*" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// specific
	w = &testResponseWriter{}
	r, _ = http.NewRequest("GET", "http://example.com/baz", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set(CORSRequestOrigin, "http://gooddomain.com")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "http://gooddomain.com" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

}
