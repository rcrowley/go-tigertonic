package tigertonic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func NotFoundHandler() http.Handler { return http.HandlerFunc(NotFound) }
