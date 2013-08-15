package tigertonic

import "net/http"

// TODO: handle other kinds of CORS headers. Early indication is that some kind
// of "interactive" way of handling certain header types will be required, so
// we may need to come up with something more robust than just dragging an
// http.Header around

const AllowOrigin string = "Access-Control-Allow-Origin"

type CORSHandler struct {
	http.Handler
	Header *http.Header
}

type CORSBuilder struct {
	*http.Header
}

func NewCORSBuilder() *CORSBuilder {
	return &CORSBuilder{&http.Header{}}
}

func (self *CORSBuilder) SetAllowedOrigin(origin string) *CORSBuilder {
	self.Header.Set(AllowOrigin, origin)
	return self
}

func (self *CORSBuilder) Build(handler http.Handler) *CORSHandler {
	return &CORSHandler{handler, self.Header}
}
