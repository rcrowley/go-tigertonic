package tigertonic

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

// Server is an http.Server with better defaults.
type Server struct {
	http.Server
	listener  net.Listener
	waitGroup sync.WaitGroup
}

// NewServer returns an http.Server with better defaults.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: http.Server{
			Addr: addr,
			Handler: &serverHandler{
				Handler: handler,
			},
			MaxHeaderBytes: 4096,
			ReadTimeout:    60e9, // These are absolute times which must be
			WriteTimeout:   60e9, // longer than the longest {up,down}load.
		},
	}
}

// NewTLSServer returns an http.Server with better defaults configured to use
// the certificate and private key files.
func NewTLSServer(
	addr, cert, key string,
	handler http.Handler,
) (*Server, error) {
	s := NewServer(addr, handler)
	return s, s.TLS(cert, key)
}

// CA overrides the certificate authority on the server's TLSConfig field.
func (s *Server) CA(ca string) error {
	certPool := x509.NewCertPool()
	buf, err := ioutil.ReadFile(ca)
	if nil != err {
		return err
	}
	certPool.AppendCertsFromPEM(buf)
	s.tlsConfig()
	s.TLSConfig.RootCAs = certPool
	return nil
}

// Close closes the listener the server is using and signals open connections
// to close at their earliest convenience.
func (s *Server) Close() error {
	return s.listener.Close()
}

// ListenAndServe calls net.Listen with s.Addr and then calls s.Serve.
func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if "" == addr {
		if nil == s.TLSConfig {
			addr = ":http"
		} else {
			addr = ":https"
		}
	}
	l, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}
	return s.Serve(l)
}

// ListenAndServeTLS calls s.TLS with the given certificate and private key
// files and then calls s.ListenAndServe.
func (s *Server) ListenAndServeTLS(cert, key string) error {
	s.TLS(cert, key)
	return s.ListenAndServe()
}

// Serve behaves like http.Server.Serve with the added option to stop the
// server gracefully with the s.Close and s.Wait methods.
func (s *Server) Serve(l net.Listener) error {
	s.listener = &Listener{
		Listener:  l,
		waitGroup: &s.waitGroup,
	}
	return s.Server.Serve(s.listener)
}

// TLS configures this server to be a TLS server using the given certificate
// and private key files.
func (s *Server) TLS(cert, key string) error {
	c, err := tls.LoadX509KeyPair(cert, key)
	if nil != err {
		return err
	}
	s.tlsConfig()
	s.TLSConfig.Certificates = []tls.Certificate{c}
	return nil
}

// Wait waits for all open connections to become closed.
func (s *Server) Wait() {
	s.waitGroup.Wait()
}

func (s *Server) tlsConfig() {
	if nil == s.TLSConfig {
		s.TLSConfig = &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
			},
		}
	}
}

type serverHandler struct {
	http.Handler
	ch <-chan struct{}
}

func (h *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// r.Header.Set("Host", r.Host) // Should I?
	r.URL.Host = r.Host
	if nil != r.TLS {
		r.URL.Scheme = "https"
	} else {
		r.URL.Scheme = "http"
	}
	s.handler.ServeHTTP(w, r)
}
