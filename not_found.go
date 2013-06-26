package tigertonic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// NotFound responds 404 to every request, possibly with a JSON body.
func NotFound(w http.ResponseWriter, r *http.Request) {
	description := fmt.Sprintf("%s %s not found", r.Method, r.URL.Path)
	if acceptJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"description": description,
			"error":       "NotFound",
		}); nil != err {
			log.Println(err)
		}
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, description)
	}
}

// NotFoundHandler responds 404 to every request, possibly with a JSON body.
func NotFoundHandler() http.Handler { return http.HandlerFunc(NotFound) }
