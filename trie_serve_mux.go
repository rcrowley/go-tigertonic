package tigertonic

import (
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// TrieServeMux is an HTTP request multiplexer that implements http.Handler
// with an API similar to http.ServeMux.  It is expanded to be sensitive to the
// HTTP method and treats URL patterns as patterns rather than simply prefixes.
//
// Components of the URL pattern surrounded by braces (for example: "{foo}")
// match any string and create an entry for the string plus the string
// surrounded by braces in the query parameters (for example: "foo" and
// "{foo}").
type TrieServeMux struct {
	methods map[string]http.Handler
	param   *string
	paths   map[string]*TrieServeMux
}

// NewTrieServeMux makes a new TrieServeMux.
func NewTrieServeMux() *TrieServeMux {
	return &TrieServeMux{
		methods: make(map[string]http.Handler),
		param:   nil,
		paths:   make(map[string]*TrieServeMux),
	}
}

// Handle registers an http.Handler for the given HTTP method and URL pattern.
func (mux *TrieServeMux) Handle(method, pattern string, handler http.Handler) {
	log.Printf("handling %s %s\n", method, pattern)
	mux.addRoute(method, strings.Split(pattern, "/")[1:], handler)
}

// HandleFunc registers a handler function for the given HTTP method and URL
// pattern.
func (mux *TrieServeMux) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(method, pattern, http.HandlerFunc(handler))
}

// ServeHTTP routes an HTTP request to the http.Handler registered for the URL
// pattern which matches the requested path.  It responds 404 if there is no
// matching URL pattern and 405 if the requested HTTP method is not allowed.
func (mux *TrieServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// r.Header.Set("Host", r.Host)
	r.URL.Host = r.Host
	if nil != r.TLS {
		r.URL.Scheme = "https"
	} else {
		r.URL.Scheme = "http"
	}
	params, handler := mux.findRoute(
		r.Method,
		strings.Split(r.URL.Path, "/")[1:],
	)
	r.URL.RawQuery = r.URL.RawQuery + "&" + params.Encode()
	handler.ServeHTTP(w, r)
}

// addRoute recursively adds a URL pattern, parsing wildcards as it goes, to
// the trie and registers an http.Handler to handle an HTTP method.
func (mux *TrieServeMux) addRoute(method string, paths []string, handler http.Handler) {
	if 0 == len(paths) {
		mux.methods[method] = handler
		return
	}
	if strings.HasPrefix(paths[0], "{") && strings.HasSuffix(paths[0], "}") {
		mux.param = &paths[0]
	}
	if _, ok := mux.paths[paths[0]]; !ok {
		mux.paths[paths[0]] = NewTrieServeMux()
	}
	mux.paths[paths[0]].addRoute(method, paths[1:], handler)
}

// findRoute recursively searches for a URL pattern in the trie, adds wildcards
// to the query parameters, and returns an http.Handler.
func (mux *TrieServeMux) findRoute(method string, paths []string) (url.Values, http.Handler) {
	if 0 == len(paths) {
		if _, ok := mux.methods[method]; !ok {
			return nil, methodNotAllowedHandler{mux}
		}
		return nil, mux.methods[method]
	}
	if _, ok := mux.paths[paths[0]]; ok {
		return mux.paths[paths[0]].findRoute(method, paths[1:])
	}
	if nil != mux.param {
		params, handler := mux.paths[*mux.param].findRoute(
			method,
			paths[1:],
		)
		if nil == params {
			params = make(url.Values)
		}
		params.Set(*mux.param, paths[0])
		params.Set(strings.Trim(*mux.param, "{}"), paths[0])
		return params, handler
	}
	return nil, http.NotFoundHandler()
}

type methodNotAllowedHandler struct {
	mux *TrieServeMux
}

func (h methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methods := []string{"OPTIONS"}
	if _, ok := h.mux.methods["GET"]; ok {
		methods = append(methods, "HEAD")
	}
	for method, _ := range h.mux.methods {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	w.Header().Set("Allow", strings.Join(methods, ", "))
	if "OPTIONS" == r.Method {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
