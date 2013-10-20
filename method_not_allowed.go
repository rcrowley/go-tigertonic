package tigertonic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

// MethodNotAllowedHandler responds 405 to every request with an Allow header
// and possibly with a JSON body.
type MethodNotAllowedHandler struct {
	mux *TrieServeMux
}

func (h MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		if method := r.Header.Get(CORSRequestMethod); method != "" {
			w.Header().Set(CORSAllowMethods, strings.Join(methods, ", "))
			if requestOrigin := r.Header.Get(CORSRequestOrigin); requestOrigin != "" {
				allowedOrigin := ""
				if cors, ok := h.mux.methods[method].(*CORSHandler); ok {
					allowedOrigin = cors.getAllowedOrigin(requestOrigin)
				}

				if allowedOrigin == "" {
					allowedOrigin = "null"
				}
				w.Header().Set(CORSAllowOrigin, allowedOrigin)
			}
			if requestHeaders := r.Header.Get(CORSRequestHeaders); requestHeaders != "" {
				allowedHeaders := ""
				if cors, ok := h.mux.methods[method].(*CORSHandler); ok {
					allowedHeaders = cors.getAllowedHeaders()
				}
				w.Header().Set(CORSAllowHeaders, allowedHeaders)
			}

		}
		if acceptJSON(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(map[string][]string{
				"allow": methods,
			}); nil != err {
				log.Println(err)
			}
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, strings.Join(methods, ", "))
		}
	} else {
		description := fmt.Sprintf(
			"only %s are allowed",
			strings.Join(methods, ", "),
		)
		if acceptJSON(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			var e string
			if SnakeCaseHTTPEquivErrors {
				e = "method_not_allowed"
			} else {
				e = "tigertonic.MethodNotAllowed"
			}
			if err := json.NewEncoder(w).Encode(map[string]string{
				"description": description,
				"error":       e,
			}); nil != err {
				log.Println(err)
			}
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, description)
		}
	}
}
