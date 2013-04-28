package main

import (
	"flag"
	"fmt"
	"github.com/rcrowley/go-tigertonic"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	cert   = flag.String("cert", "", "certificate pathname")
	key    = flag.String("key", "", "private key pathname")
	listen = flag.String("listen", "127.0.0.1:8000", "listen address")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: example [-listen=<listen>]")
		flag.PrintDefaults()
	}
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	flag.Parse()
	mux := tigertonic.NewTrieServeMux()
	mux.Handle("POST", "/stuff", tigertonic.Marshaled(create))
	mux.Handle("GET", "/stuff/{id}", tigertonic.Marshaled(get))
	mux.Handle("POST", "/stuff/{id}", tigertonic.Marshaled(update))
	server := &http.Server{
		Addr:           *listen,
		Handler:        tigertonic.Logged(tigertonic.Timed(mux)),
		MaxHeaderBytes: 4096,
		ReadTimeout:    1e9,
		WriteTimeout:   1e9,
	}
	if "" != *cert && "" != *key {
		server.ListenAndServeTLS(*cert, *key)
	} else {
		server.ListenAndServe()
	}
}

// POST /stuff
func create(u *url.URL, h http.Header, rq *MyRequest) (int, http.Header, *MyResponse, error) {
	return http.StatusCreated, http.Header{
		"Content-Location": {fmt.Sprintf(
			"%s://%s/stuff/%s",
			u.Scheme,
			u.Host,
			rq.ID,
		)},
	}, &MyResponse{rq.ID, rq.Stuff}, nil
}

// GET /stuff/{id}
func get(u *url.URL, h http.Header, rq *MyRequest) (int, http.Header, *MyResponse, error) {
	return http.StatusOK, nil, &MyResponse{u.Query().Get("id"), "STUFF"}, nil
}

// POST /stuff/{id}
func update(u *url.URL, h http.Header, rq *MyRequest) (int, http.Header, *MyResponse, error) {
	return http.StatusAccepted, nil, &MyResponse{u.Query().Get("id"), "STUFF"}, nil
}
