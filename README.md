Tiger Tonic
===========

A Go framework for building JSON web services.  Inspired by [Dropwizard](http://dropwizard.codahale.com).

Like the Go language itself, Tiger Tonic strives to keep features orthogonal.

`TrieServeMux`
--------------

HTTP routing in the Go standard library is pretty anemic.  Enter `TrieServeMux`.  It accepts an HTTP method, a URL pattern, and an `http.Handler` or an `http.HandlerFunc`.  Components in the URL pattern wrapped in `{` and `}` are wildcards: their values are added to the URL as <code>u.Query().Get("<em>name</em>")</code>.

`Marshaled`
-----------

Wrap a function in `Marshaled` to turn it into an `http.Handler`.  The function signature must be something like this or `Marshaled` will panic:

```go
func myHandler(*url.URL, http.Header, *MyRequest) (int, http.Header, *MyResponse)
```

Request bodies will be unmarshaled into a `MyRequest` struct and response bodies will be marshaled from `MyResponse` structs.

`Logged`
--------

Wrap an `http.Handler` in `Logged` to have the request and response headers and bodies logged to standard error.

`Timed`
-------

Wrap an `http.Handler` in `Timed` to have the request timed with [`go-metrics`](https://github.com/rcrowley/go-metrics).

Usage
-----

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
listener, _ := net.ListenTCP("tcp", laddr)
mux := NewTrieServeMux()
mux.Handle("GET", "/stuff", Marshaled(myHandler))
server := &http.Server{
    Addr:           laddr.String(),
    Handler:        Logged(Timed(mux)),
    MaxHeaderBytes: 4096,
    ReadTimeout:    1e9,
    WriteTimeout:   1e9,
}
server.Serve(listener)
```

Ready for more?  See the full [example](https://github.com/rcrowley/go-tigertonic/tree/master/example).

WTF?
----

Dropwizard was named by <http://gunshowcomic.com/316> so Tiger Tonic was named by <http://gunshowcomic.com/338>.
