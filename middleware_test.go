package tigertonic

import (
	"errors"
	"net/http"
	"testing"
)

func TestFirst1(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(NotFoundHandler()).ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

func TestFirst2(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(noopHandler{}, NotFoundHandler()).ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

func TestFirst3(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(noopHandler{}, noopHandler{}, NotFoundHandler()).ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

func TestFirst4(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	First(NotFoundHandler(), &fatalHandler{t}).ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

func TestIfFalse(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	If(func(r *http.Request) error {
		return Unauthorized{errors.New("Unauthorized")}
	}, NotFoundHandler()).ServeHTTP(w, r)
	if http.StatusUnauthorized != w.Status {
		t.Fatal(w.Status)
	}
}

func TestIfTrue(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	If(func(r *http.Request) error {
		return nil
	}, NotFoundHandler()).ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

type fatalHandler struct {
	t *testing.T
}

func (fh *fatalHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	fh.t.Fatal("fatalHandler")
}

type noopHandler struct{}

func (noopHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
