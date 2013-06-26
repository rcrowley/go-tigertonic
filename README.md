Tiger Tonic
===========

A Go framework for building JSON web services inspired by [Dropwizard](http://dropwizard.codahale.com).  If HTML is your game, this will hurt a little.

Like the Go language itself, Tiger Tonic strives to keep features orthogonal.  It defers what it can to the Go standard library and a few other packages.

`tigertonic.TrieServeMux`
-------------------------

HTTP routing in the Go standard library is pretty anemic.  Enter `tigertonic.TrieServeMux`.  It accepts an HTTP method, a URL pattern, and an `http.Handler` or an `http.HandlerFunc`.  Components in the URL pattern wrapped in `{` and `}` are wildcards: their values are added to the URL as <code>u.Query().Get("<em>name</em>")</code>.

`tigertonic.HostServeMux`
-------------------------

Use `tigertonic.HostServeMux` to serve multiple domain names from the same `net.Listener`.

`tigertonic.Marshaled`
----------------------

Wrap a function in `tigertonic.Marshaled` to turn it into an `http.Handler`.  The function signature must be something like this or `tigertonic.Marshaled` will panic:

```go
func myHandler(*url.URL, http.Header, *MyRequest) (int, http.Header, *MyResponse)
```

Request bodies will be unmarshaled into a `MyRequest` struct and response bodies will be marshaled from `MyResponse` structs.

`tigertonic.Logged`
-------------------

Wrap an `http.Handler` in `tigertonic.Logged` to have the request and response headers and bodies logged to standard error.  The second argument is an optional `func(string) string` called to as requests and responses are logged to give the caller the opportunity to redact sensitive information from log entries.

`tigertonic.Counted` and `tigertonic.Timed`
-------------------------------------------

Wrap an `http.Handler` in `Counted` or `Timed` to have the request counted or timed with [`go-metrics`](https://github.com/rcrowley/go-metrics).

Usage
-----

Install dependencies:

```sh
sh bootstrap.sh
```

Then define your service.  The working [example](https://github.com/rcrowley/go-tigertonic/tree/master/example) may be a more convenient place to start.

Requests that have bodies have types.  JSON is deserialized by adding `tigertonic.Marshaled` to your routes.

```go
type MyRequest struct {
	ID     string      `json:"id"`
	Stuff  interface{} `json:"stuff"`
}
```

Responses, too, have types.  JSON is serialized by adding `tigertonic.Marshaled` to your routes.

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

Wire it all up in `main.main`!

```go
mux := NewTrieServeMux()
mux.Handle("GET", "/stuff", tigertonic.Marshaled(tigertonic.Timed(myHandler, "myHandler", nil)))
tigertonic.NewServer(":8000", tigertonic.Logged(mux, nil)).ListenAndServe()
```

Ready for more?  See the full [example](https://github.com/rcrowley/go-tigertonic/tree/master/example).

WTF?
----

Dropwizard was named after <http://gunshowcomic.com/316> so Tiger Tonic was named after <http://gunshowcomic.com/338>.

If Tiger Tonic isn't your cup of tea, perhaps one of these fine tools suits you:

* <https://github.com/bmizerany/pat>
* <https://github.com/hoisie/web.go>
* <http://www.gorillatoolkit.org>
* <http://www.gorillatoolkit.org>
