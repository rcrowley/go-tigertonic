package tigertonic

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
)

// HTTPBasicAuth returns an http.Handler that conditionally calls another
// http.Handler if the request includes and Authorization header with a
// username and password that appear in the map of credentials.  Otherwise,
// respond 401 Unauthorized.
//
// TODO Respond with an appropriate WWW-Authenticate header.
func HTTPBasicAuth(credentials map[string]string, h http.Handler) FirstHandler {
	return If(func(r *http.Request) error {
		auth := r.Header.Get("Authorization")
		if 6 > len(auth) || "Basic " != auth[:6] {
			return Unauthorized{errors.New("no HTTP Basic auth specified")}
		}
		buf, err := base64.StdEncoding.DecodeString(auth[6:])
		if nil != err {
			return Unauthorized{err}
		}
		i := bytes.IndexByte(buf, ':')
		if -1 == i {
			return Unauthorized{errors.New("malformed HTTP Basic auth specified")}
		}
		if password, ok := credentials[string(buf[:i])]; !ok || password != string(buf[i+1:]) {
			return Unauthorized{errors.New("unauthorized")}
		}
		return nil
	}, h)
}
