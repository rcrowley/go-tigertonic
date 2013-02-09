package tigertonic

import (
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type TrieServeMux struct {
	methods map[string]http.Handler
	paths   map[string]*TrieServeMux
	param   *string
}

func NewTrieServeMux() *TrieServeMux {
	return &TrieServeMux{
		make(map[string]http.Handler),
		make(map[string]*TrieServeMux),
		nil,
	}
}

func (mux *TrieServeMux) Handle(method, pattern string, handler http.Handler) {
	log.Printf("handling %s %s\n", method, pattern)
	mux.addRoute(method, strings.Split(pattern, "/")[1:], handler)
}

func (mux *TrieServeMux) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(method, pattern, http.HandlerFunc(handler))
}

func (mux *TrieServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, handler := mux.findRoute(
		r.Method,
		strings.Split(r.URL.Path, "/")[1:],
	)
	r.URL.RawQuery = r.URL.RawQuery + "&" + params.Encode()
	handler.ServeHTTP(w, r)
}

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
