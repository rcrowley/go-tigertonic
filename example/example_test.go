package main

import (
	"github.com/rcrowley/go-tigertonic/mocking"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	s, h, rs, err := create(
		mocking.URL(hMux, "POST", "http://example.com/1.0/stuff"),
		mocking.Header(nil),
		&MyRequest{"ID", "STUFF"},
	)
	if nil != err {
		t.Fatal(err)
	}
	if http.StatusCreated != s {
		t.Fatal(s)
	}
	if "http://example.com/1.0/stuff/ID" != h.Get("Content-Location") {
		t.Fatal(h)
	}
	if "ID" != rs.ID || "STUFF" != rs.Stuff {
		t.Fatal(rs)
	}
}

func TestGet(t *testing.T) {
	s, _, rs, err := get(
		mocking.URL(hMux, "GET", "http://example.com/1.0/stuff/ID"),
		mocking.Header(nil),
		nil,
	)
	if nil != err {
		t.Fatal(err)
	}
	if http.StatusOK != s {
		t.Fatal(s)
	}
	if "ID" != rs.ID || "STUFF" != rs.Stuff {
		t.Fatal(rs)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	defer func() { recover() }()
	mocking.URL(hMux, "PUT", "http://example.com/1.0/stuff")
	t.Fail()
}

func TestNotFound(t *testing.T) {
	defer func() { recover() }()
	mocking.URL(hMux, "GET", "http://example.com/1.0/things")
	t.Fail()
}

func TestUpdate(t *testing.T) {
	s, _, rs, err := update(
		mocking.URL(hMux, "POST", "http://example.com/1.0/stuff/ID"),
		mocking.Header(nil),
		&MyRequest{"ID", "STUFF"},
	)
	if nil != err {
		t.Fatal(err)
	}
	if http.StatusAccepted != s {
		t.Fatal(s)
	}
	if "ID" != rs.ID || "STUFF" != rs.Stuff {
		t.Fatal(rs)
	}
}
