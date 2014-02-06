package tigertonic

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

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

// ClientCA configures the CA pool for verifying client side certificates
func (s *Server) ClientCA(ca string) error {
	certPool := x509.NewCertPool()
	buf, err := ioutil.ReadFile(ca)
	if nil != err {
		return err
	}

	certPool.AppendCertsFromPEM(buf)
	s.tlsConfig()
	s.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	s.TLSConfig.ClientCAs = certPool
	return nil
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
