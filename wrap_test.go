package tigertonic

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

type TestWrapper struct {
	http.ResponseWriter
}

type TestPostProcessor struct {
	http.Handler
	Processed int64
}

func (self *TestPostProcessor) WrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &TestWrapper{w}
}

func (self *TestPostProcessor) Process(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(*TestWrapper); ok {
		self.Processed++
	} else {
		self.Processed--
	}
}

func TestWrap(t *testing.T) {
	var postProcessor PostProcessor = &TestPostProcessor{}
	if _, ok := postProcessor.(ResponseWriterWrapper); !ok {
		t.Fatal("cant cast")
	}
	h := Processed(Marshaled(func(u *url.URL, h http.Header, rq *testRequest) (int, http.Header, *testResponse, error) {
		return http.StatusOK, nil, &testResponse{"bar"}, nil
	}), postProcessor)
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", bytes.NewBufferString(`{"foo":"bar"}`))
	h.ServeHTTP(w, r)
	if -1 == postProcessor.(*TestPostProcessor).Processed {
		t.Fatal("Post processing occurred, but wrapping did not")
	} else if 0 == postProcessor.(*TestPostProcessor).Processed {
		t.Fatal("Post processing did not occur")
	}
}
