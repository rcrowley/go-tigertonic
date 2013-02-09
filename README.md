Tiger Tonic
===========

A Go framework for building JSON web services.  Inspired by [Dropwizard](http://dropwizard.codahale.com).

Synopsis
--------

Description
-----------

Examples
--------

Requests that have bodies have types.  JSON is deserialized by adding `Marshaled` to your routes.

```go
type MyRequest struct {
	ID     string      `json:"id"`
	Stuff  interface{} `json:"stuff"`
}
```

Responses, too, have types.  JSON is serialized by adding `Marshaled` to your routes.

```go
type MyResponse struct {
	ID     string      `json:"id"`
	Stuff  interface{} `json:"stuff"`
}
```

Routes are just functions with a particular signature.  You control the request and response types.

```go
func myHandler(u *url.URL, h http.Header, *MyRequest) (int, http.Header, *MyResponse, error) {
    return http.StatusOK, nil, &MyResponse{"ID", "STUFF"}, nil
}
```

Wire it all up in `main`!

```go
laddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8000")
mux := NewTrieServeMux()
mux.Handle("GET", "/stuff", Marshaled(myHandler))
server := &http.Server{
    Addr:           laddr.String(),
    Handler:        Logged(Timed(mux)),
    MaxHeaderBytes: 4096,
    ReadTimeout:    1e9,
    WriteTimeout:   1e9,
}
server.Serve(laddr)
```

(Don't put this into production.  See the full [example](https://github.com/rcrowley/go-tigertonic/tree/master/example), for error handling.)

WTF?
----

Dropwizard was named by <http://gunshowcomic.com/316> so Tiger Tonic was named by <http://gunshowcomic.com/338>.
