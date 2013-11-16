package tigertonic

import "testing"

func TestServerCATLS(t *testing.T) {
	s := NewServer("", NotFoundHandler{})
	s.CA("test.crt")
	s.TLS("test.crt", "test.key")
	if nil == s.TLSConfig.Certificates || 1 != len(s.TLSConfig.Certificates) {
		t.Fatal("no Certificates")
	}
	if nil == s.TLSConfig.RootCAs || 1 != len(s.TLSConfig.RootCAs.Subjects()) {
		t.Fatal("no RootCAs")
	}
}
