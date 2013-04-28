package tigertonic

import (
	"net/http"
	"testing"
)

func TestNotFound(t *testing.T) {
	mux := NewTrieServeMux()
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	mux.ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := NewTrieServeMux()
	mux.HandleFunc("POST", "/", func(w http.ResponseWriter, r *http.Request) {})
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	mux.ServeHTTP(w, r)
	if http.StatusMethodNotAllowed != w.Status {
		t.Fatal(w.Status)
	}
}

func TestOPTIONS(t *testing.T) {
	mux := NewTrieServeMux()
	mux.HandleFunc("GET", "/foo", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST", "/bar", func(w http.ResponseWriter, r *http.Request) {})
	w := &testResponseWriter{}
	r, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "GET, HEAD, OPTIONS" != w.Header().Get("Allow") {
		t.Fatal(w.Header().Get("Allow"))
	}
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/bar", nil)
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.Status {
		t.Fatal(w.Status)
	}
	if "OPTIONS, POST" != w.Header().Get("Allow") {
		t.Fatal(w.Header().Get("Allow"))
	}
}

func TestRoot(t *testing.T) {
	mux := NewTrieServeMux()
	mux.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	mux.ServeHTTP(w, r)
	if http.StatusNoContent != w.Status {
		t.Fatal(w.Status)
	}
}

func TestRecurse(t *testing.T) {
	mux := NewTrieServeMux()
	mux.HandleFunc("GET", "/foo/bar/baz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo/bar/baz", nil)
	mux.ServeHTTP(w, r)
	if http.StatusNoContent != w.Status {
		t.Fatal(w.Status)
	}
}

func TestParams(t *testing.T) {
	mux := NewTrieServeMux()
	mux.HandleFunc("GET", "/{foo}/{bar}", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if "bar" != q.Get("foo") || "foo" != q.Get("bar") {
			t.Fatal(q.Get("foo"), q.Get("bar"))
		}
		if "bar" != q.Get("{foo}") || "foo" != q.Get("{bar}") {
			t.Fatal(q.Get("{foo}"), q.Get("{bar}"))
		}
		if "quux" != q.Get("baz") {
			t.Fatal(q.Get("quux"))
		}
		w.WriteHeader(http.StatusNoContent)
	})
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/bar/foo?baz=quux", nil)
	mux.ServeHTTP(w, r)
	if http.StatusNoContent != w.Status {
		t.Fatal(w.Status)
	}
}
