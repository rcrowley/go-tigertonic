package tigertonic

import (
	"log"
	"net/http"
)

// HostServeMux is an HTTP request multiplexer that implements http.Handler
// with an API similar to http.ServeMux.  It is only sensitive to the hostname
// and doesn't even look at the rest of the request.
type HostServeMux map[string]http.Handler

// NewHostServeMux makes a new HostServeMux.
func NewHostServeMux() HostServeMux {
	return make(HostServeMux)
}

// Handle registers an http.Handler for the given hostname.
func (mux HostServeMux) Handle(hostname string, handler http.Handler) {
	log.Printf("handling %s\n", hostname)
	mux[hostname] = handler
}

// HandleFunc registers a handler function for the given hostname.
func (mux HostServeMux) HandleFunc(hostname string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(hostname, http.HandlerFunc(handler))
}

// ServeHTTP routes an HTTP request to the http.Handler registered for the
// requested hostname.  It responds 404 if the hostname is not registered.
func (mux HostServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := mux[r.Host]
	if !ok {
		handler = NotFoundHandler()
	}
	handler.ServeHTTP(w, r)
}
