package tigertonic

import (
	"net/http"
)

// An abstraction for wrapping request/response pairs to make it easier to
// write additional pre/post processors

type ResponseWriterWrapper interface {
	WrapResponseWriter(http.ResponseWriter) http.ResponseWriter
}

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
