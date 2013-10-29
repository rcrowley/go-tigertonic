package tigertonic

import "testing"

func TestServerCATLS(t *testing.T) {
	s := NewServer("", NotFoundHandler{})
	s.CA("/etc/ssl/certs/betable-ca.pem")
	s.TLS(
		"/etc/ssl/certs/betable-internal.crt",
		"/etc/ssl/private/betable-internal.key",
	)
	if nil == s.TLSConfig.Certificates || 1 != len(s.TLSConfig.Certificates) {
		t.Fatal("no Certificates")
	}
	if nil == s.TLSConfig.RootCAs || 1 != len(s.TLSConfig.RootCAs.Subjects()) {
		t.Fatal("no RootCAs")
	}
}
