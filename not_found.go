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
		ResponseErrorWriter.WriteJSONError(w, NotFound{Err: errors.New(description)})
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, description)
	}
}
