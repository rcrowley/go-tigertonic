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
		Handler:        Logged(handler),
		MaxHeaderBytes: 4096,
		ReadTimeout:    1e9,
		WriteTimeout:   1e9,
	}}
}
