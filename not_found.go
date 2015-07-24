package tigertonic

import (
	"errors"
	"fmt"
	"net/http"
)

// NotFoundHandler responds 404 to every request, possibly with a JSON body.
type NotFoundHandler struct{}

func (NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	description := fmt.Sprintf("%s %s not found", r.Method, r.URL.Path)
	if acceptJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := NotFound{Err: errors.New(description)}
		ResponseErrorWriter.WriteJSONError(w, err)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, description)
	}
}
