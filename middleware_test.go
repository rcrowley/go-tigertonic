package tigertonic

import (
	"net/http"
	"testing"
)

func TestFirst1(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(NotFoundHandler()).ServeHTTP(w, r)
	if 404 != w.Status {
		t.Fatal(w.Status)
	}
}

func TestFirst2(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(noopHandler{}, NotFoundHandler()).ServeHTTP(w, r)
	if 404 != w.Status {
		t.Fatal(w.Status)
	}
}

func TestFirst3(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(noopHandler{}, noopHandler{}, NotFoundHandler()).ServeHTTP(w, r)
	if 404 != w.Status {
		t.Fatal(w.Status)
	}
}

func TestFirst4(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(NotFoundHandler(), &fatalHandler{t}).ServeHTTP(w, r)
	if 404 != w.Status {
		t.Fatal(w.Status)
	}
}

type fatalHandler struct{
	t *testing.T
}

func (fh *fatalHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	fh.t.Fatal("fatalHandler")
}

type noopHandler struct{}

func (noopHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
