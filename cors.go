package tigertonic

import (
	"net/http"
	"strings"
)

// BUG(mihasya): currently only handles Origin and header-related CORS headers.
// Early indication is that logic will be required for handling certain header
// types (the credential-related ones in particular) so we may need to come up
// with something more robust than just dragging an http.Header around

const CORSRequestOrigin string = "Origin"
const CORSRequestMethod string = "Access-Control-Request-Method"
const CORSRequestHeaders string = "Access-Control-Request-Headers"

const CORSAllowOrigin string = "Access-Control-Allow-Origin"
const CORSAllowMethods string = "Access-Control-Allow-Methods"
const CORSAllowHeaders string = "Access-Control-Allow-Headers"

// CORSHandler wraps an http.Handler while correctly handling CORS related
// functionality, such as Origin headers. It also allows tigertonic core to
// correctly respond to OPTIONS headers for CORS-enabled endpoints
type CORSHandler struct {
	http.Handler
	Header http.Header
}

func (self *CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if requestOrigin := r.Header.Get("Origin"); requestOrigin != "" {
		w.Header().Set(CORSAllowOrigin, self.getAllowedOrigin(requestOrigin))
	}
	self.Handler.ServeHTTP(w, r)
}

func (self *CORSHandler) getAllowedOrigin(requestOrigin string) string {
	if self.Header.Get(CORSAllowOrigin) == "*" {
		return "*"
	} else if self.Header.Get(CORSAllowOrigin) == requestOrigin {
		return requestOrigin
	}
	return "null"
}

func (self *CORSHandler) getAllowedHeaders() string {
	return self.Header.Get(CORSAllowHeaders)
}

// CRSBuilder facilitates the application of the same set of CORS rules to a
// number of endpoints. One would use CORSBuilder.Build() the same way one
// might wrap a handler in a call to Timed() or Logged().
type CORSBuilder struct {
	http.Header
	allowedHeaders []string
}

func NewCORSBuilder() *CORSBuilder {
	return &CORSBuilder{http.Header{}, []string{}}
}

// SetAllowedOrigin sets the domain for which cross-origin requests are allowed
func (self *CORSBuilder) SetAllowedOrigin(origin string) *CORSBuilder {
	self.Header.Set(CORSAllowOrigin, origin)
	return self
}

func (self *CORSBuilder) AddAllowedHeaders(headers ...string) *CORSBuilder {
	self.allowedHeaders = append(self.allowedHeaders, headers...)
	return self
}

func (self *CORSBuilder) Build(handler http.Handler) *CORSHandler {
	self.Header.Set(CORSAllowHeaders, strings.Join(self.allowedHeaders, ", "))
	return &CORSHandler{handler, self.Header}
}
