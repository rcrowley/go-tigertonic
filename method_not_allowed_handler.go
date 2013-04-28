package tigertonic

import (
	"net/http"
	"sort"
	"strings"
)

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
