package tigertonic

import (
	"net/http"
)

// An abstraction for wrapping request/response pairs to make it easier to
// write additional pre/post processors

// ResponseWriterWrapper can be optionally implemented by a PostProcessor in
// order to wrap or replace http.ResponseWriter before it's passed to an
// http.Handler
type ResponseWriterWrapper interface {
	WrapResponseWriter(http.ResponseWriter) http.ResponseWriter
}

// PostProcessor is similar to http.Handler but it is expected to not actually
// interact with the request, merely perform additional accounting is needed
type PostProcessor interface {
	Process(http.ResponseWriter, *http.Request)
}

func Processed(h http.Handler, p PostProcessor) http.Handler {
	return &ProcessedHandler{h, p}
}

type ProcessedHandler struct {
	h http.Handler
	p PostProcessor
}

func (self *ProcessedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var wrapped http.ResponseWriter
	if wrapper, ok := self.p.(ResponseWriterWrapper); ok {
		wrapped = wrapper.WrapResponseWriter(w)
	} else {
		wrapped = w
	}
	self.h.ServeHTTP(wrapped, r)
	self.p.Process(wrapped, r)
}
