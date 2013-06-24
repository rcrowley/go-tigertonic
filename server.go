package tigertonic

import (
	"net/http"
)

type Server struct {
	http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{http.Server{
		Addr:           addr,
		Handler:        &server{handler},
		MaxHeaderBytes: 4096,
		ReadTimeout:    1e9,
		WriteTimeout:   1e9,
	}}
}

type server struct {
	handler http.Handler
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// r.Header.Set("Host", r.Host) // Should I?
	r.URL.Host = r.Host
	if nil != r.TLS {
		r.URL.Scheme = "https"
	} else {
		r.URL.Scheme = "http"
	}
	s.handler.ServeHTTP(w, r)
}
