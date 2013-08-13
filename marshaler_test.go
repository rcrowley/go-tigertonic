package tigertonic

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestMarshaledPanicNumIn(t *testing.T) {
	testMarshaledPanic(func() {}, t)
	testMarshaledPanic(func(u int) {}, t)
	testMarshaledPanic(func(u, h int) {}, t)
	testMarshaledPanic(func(u, h, rq, foo int) {}, t)
}

func TestMarshaledPanicIn0(t *testing.T) {
	testMarshaledPanic(func(u, h, rq int) {}, t)
}

func TestMarshaledPanicIn1(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h, rq int) {}, t)
}

func TestMarshaledPanicIn2(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq int) {}, t)
}

func TestMarshaledPanicIn3(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) {}, t)
}

func TestMarshaledPanicNumOut(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) {}, t)
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) int {
		return 0
	}, t)
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (int, int) {
		return 0, 0
	}, t)
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (int, int, int) {
		return 0, 0, 0
	}, t)
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (int, int, int, int, int) {
		return 0, 0, 0, 0, 0
	}, t)
}

func TestMarshaledPanicOut0(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (string, int, int, int) {
		return "", 0, 0, 0
	}, t)
}

func TestMarshaledPanicOut1(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (int, int, int, int) {
		return 0, 0, 0, 0
	}, t)
}

func TestMarshaledPanicOut2(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, int, int) {
		return 0, http.Header{}, 0, 0
	}, t)
}

func TestMarshaledPanicOut3(t *testing.T) {
	testMarshaledPanic(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, int) {
		return 0, http.Header{}, nil, 0
	}, t)
}

func TestNotAcceptable(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusNoContent, nil, nil, nil
	}).ServeHTTP(w, r)
	if http.StatusNotAcceptable != w.Status {
		t.Fatal(w.Status)
	}
}

func TestUnsupportedMediaType(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusNoContent, nil, nil, nil
	}).ServeHTTP(w, r)
	if http.StatusUnsupportedMediaType != w.Status {
		t.Fatal(w.Status)
	}
}

func TestBadRequest(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(""))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusNoContent, nil, nil, nil
	}).ServeHTTP(w, r)
	if http.StatusBadRequest != w.Status {
		t.Fatal(w.Status)
	}
	if "{\"description\":\"EOF\",\"error\":\"error\"}\n" != w.Body.String() {
		t.Fatal(w.Body.String())
	}
}

func TestBadRequestSyntaxError(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString("}"))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusNoContent, nil, nil, nil
	}).ServeHTTP(w, r)
	if http.StatusBadRequest != w.Status {
		t.Fatal(w.Status)
	}
	if "{\"description\":\"invalid character '}' looking for beginning of value\",\"error\":\"json.SyntaxError\"}\n" != w.Body.String() {
		t.Fatal(w.Body.String())
	}
}

func TestInternalServerError(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return 0, nil, nil, errors.New("foo")
	}).ServeHTTP(w, r)
	if http.StatusInternalServerError != w.Status {
		t.Fatal(w.Status)
	}
	if "{\"description\":\"foo\",\"error\":\"error\"}\n" != w.Body.String() {
		t.Fatal(w.Body.String())
	}
}

func TestHTTPEquivError(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return 0, nil, nil, ServiceUnavailable{errors.New("foo")}
	}).ServeHTTP(w, r)
	if http.StatusServiceUnavailable != w.Status {
		t.Fatal(w.Status)
	}
	if "{\"description\":\"foo\",\"error\":\"tigertonic.ServiceUnavailable\"}\n" != w.Body.String() {
		t.Fatal(w.Body.String())
	}
}

func TestSnakeCaseHTTPEquivError(t *testing.T) {
	SnakeCaseHTTPEquivErrors = true
	defer func() { SnakeCaseHTTPEquivErrors = false }()
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return 0, nil, nil, ServiceUnavailable{errors.New("foo")}
	}).ServeHTTP(w, r)
	if http.StatusServiceUnavailable != w.Status {
		t.Fatal(w.Status)
	}
	if "{\"description\":\"foo\",\"error\":\"service_unavailable\"}\n" != w.Body.String() {
		t.Fatal(w.Body.String())
	}
}

func TestNoContent(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusNoContent, nil, nil, nil
	}).ServeHTTP(w, r)
	if http.StatusNoContent != w.Status {
		t.Fatal(w.Status)
	}
}

func TestHeader(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusNoContent, map[string][]string{
			"Foo": {"bar"},
		}, nil, nil
	}).ServeHTTP(w, r)
	if "bar" != w.Header().Get("Foo") {
		t.Fatal(w.Header().Get("Foo"))
	}
}

func TestBody(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString("{\"foo\":\"bar\"}"))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		if "bar" != rq.Foo {
			t.Fatal(rq.Foo)
		}
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}).ServeHTTP(w, r)
	if "{\"foo\":\"bar\"}\n" != w.Body.String() {
		t.Fatal(w.Body.String())
	}
}

func Test500OnMisconfiguredPost(t *testing.T) {
	w := &testResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString("anything"))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	Marshaled(func(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}).ServeHTTP(w, r)
	if http.StatusInternalServerError != w.Status {
		t.Fatalf("Server did not 500 when trying to handle a POST to a handler with interface{} as the request type")
	}
}

func testMarshaledPanic(i interface{}, t *testing.T) {
	defer func() {
		err := recover()
		if nil == err {
			t.Fail()
		}
		if _, ok := err.(MarshalerError); !ok {
			t.Error(err)
		}
	}()
	Marshaled(i)
}

type testRequest struct {
	Foo string `json:"foo"`
}

type testResponse struct {
	Foo string `json:"foo"`
}
