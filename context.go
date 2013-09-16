package tigertonic

import (
	"net/http"
	"reflect"
	"sync"
)

var (
	contexts map[*http.Request]interface{}
	mutex sync.Mutex
)

func Context(r *http.Request) interface{} {
	mutex.Lock()
	defer mutex.Unlock()
	return contexts[r]
}

type ContextHandler struct {
	handler http.Handler
	t reflect.Type
}

func WithContext(handler http.Handler, i interface{}) *ContextHandler {
	return &ContextHandler{
		handler: handler,
		t: reflect.TypeOf(i),
	}
}

func (ch *ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	contexts[r] = reflect.New(ch.t).Interface()
	mutex.Unlock()
	ch.handler.ServeHTTP(w, r)
	mutex.Lock()
	delete(contexts, r)
	mutex.Unlock()
}

func init() {
	contexts = make(map[*http.Request]interface{})
}
