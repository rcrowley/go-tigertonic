package tigertonic

import "testing"

func TestServerCATLS(t *testing.T) {
	s, err := NewTLSServer("", "test.crt", "test.key", NotFoundHandler{})
	if nil != err {
		t.Fatal(err)
	}
	s.CA("test.crt")
	if nil == s.TLSConfig.Certificates || 1 != len(s.TLSConfig.Certificates) {
		t.Fatal("no Certificates")
	}
	if nil == s.TLSConfig.RootCAs || 1 != len(s.TLSConfig.RootCAs.Subjects()) {
		t.Fatal("no RootCAs")
	}
}
