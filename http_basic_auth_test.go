package tigertonic

import (
	"encoding/base64"
	"net/http"
	"testing"
)

func TestHTTPBasicAuthAuthorized(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	r.SetBasicAuth("username", "password")
	HTTPBasicAuth(
		map[string]string{"username": "password"},
		"Tiger Tonic",
		NotFoundHandler{},
	).ServeHTTP(w, r)
	if http.StatusNotFound != w.Status {
		t.Fatal(w.Status)
	}
}

func TestHTTPBasicAuthBase64Error(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	r.Header.Set("Authorization", "Basic not-base64")
	HTTPBasicAuth(
		map[string]string{"username": "password"},
		"Tiger Tonic",
		NotFoundHandler{},
	).ServeHTTP(w, r)
	if http.StatusUnauthorized != w.Status {
		t.Fatal(w.Status)
	}
}

func TestHTTPBasicAuthMalformed(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	r.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("username")))
	HTTPBasicAuth(
		map[string]string{"username": "password"},
		"Tiger Tonic",
		NotFoundHandler{},
	).ServeHTTP(w, r)
	if http.StatusUnauthorized != w.Status {
		t.Fatal(w.Status)
	}
}

func TestHTTPBasicAuthUnauthorized(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	r.SetBasicAuth("username", "wrong-password")
	HTTPBasicAuth(
		map[string]string{"username": "password"},
		"Tiger Tonic",
		NotFoundHandler{},
	).ServeHTTP(w, r)
	if http.StatusUnauthorized != w.Status {
		t.Fatal(w.Status)
	}
}

func TestHTTPBasicAuthUnspecified(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	HTTPBasicAuth(
		map[string]string{"username": "password"},
		"Tiger Tonic",
		NotFoundHandler{},
	).ServeHTTP(w, r)
	if http.StatusUnauthorized != w.Status {
		t.Fatal(w.Status)
	}
}
