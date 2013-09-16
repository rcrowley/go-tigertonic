package tigertonic

import "net/http"

// Server is an http.Server with better defaults.
type Server struct {
	http.Server
}

// NewServer returns an http.Server with better defaults.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{http.Server{
		Addr:           addr,
		Handler:        &server{handler},
		MaxHeaderBytes: 4096,
		ReadTimeout:    60e9, // These are absolute times which must be
		WriteTimeout:   60e9, // longer than the longest {up,down}load.
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
