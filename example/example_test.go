package main

import (
	"github.com/rcrowley/go-tigertonic/mocking"
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
	if 201 != s {
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
	if 200 != s {
		t.Fatal(s)
	}
	if "ID" != rs.ID || "STUFF" != rs.Stuff {
		t.Fatal(rs)
	}
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
	if 202 != s {
		t.Fatal(s)
	}
	if "ID" != rs.ID || "STUFF" != rs.Stuff {
		t.Fatal(rs)
	}
}
