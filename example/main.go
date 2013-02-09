package main

import (
	"flag"
	"fmt"
	"github.com/rcrowley/go-tigertonic"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

var listen = flag.String("listen", "127.0.0.1:8000", "listen address")

func main() {

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: example [-listen=<listen>]")
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	laddr, err := net.ResolveTCPAddr("tcp", *listen)
	if nil != err {
		log.Fatalln(err)
	}

	mux := tigertonic.NewTrieServeMux()
	mux.Handle("POST", "/stuff", tigertonic.Marshaled(create))
	mux.Handle("GET", "/stuff/{id}", tigertonic.Marshaled(get))
	mux.Handle("POST", "/stuff/{id}", tigertonic.Marshaled(update))

	server := &http.Server{
		Addr:           laddr.String(),
		Handler:        tigertonic.Logged(tigertonic.Timed(mux)),
		MaxHeaderBytes: 4096,
		ReadTimeout:    1e9,
		WriteTimeout:   1e9,
	}
	l, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	if err := server.Serve(l); nil != err {
		log.Fatalln(err)
	}

}

// POST /stuff
func create(u *url.URL, h http.Header, rq *MyRequest) (int, http.Header, *MyResponse, error) {
	var err error = nil // TODO
	if nil != err {
		return http.StatusInternalServerError, nil, nil, err
	}
	return http.StatusCreated, map[string][]string{
		"Content-Location": {fmt.Sprintf(
			"%s://%s/stuff/%s",
			u.Scheme,
			h.Get("Host"),
			"ID",
		)},
	}, &MyResponse{"ID", "STUFF"}, nil
}

// GET /stuff/{id}
func get(u *url.URL, h http.Header, rq *MyRequest) (int, http.Header, *MyResponse, error) {
	id := u.Query().Get("id")
	var err error = nil // TODO
	if nil != err {
		return http.StatusInternalServerError, nil, nil, err
	}
	return http.StatusOK, nil, &MyResponse{id, "STUFF"}, nil
}

// POST /stuff/{id}
func update(u *url.URL, h http.Header, rq *MyRequest) (int, http.Header, *MyResponse, error) {
	id := u.Query().Get("id")
	var err error = nil // TODO
	if nil != err {
		return http.StatusInternalServerError, nil, nil, err
	}
	return http.StatusOK, nil, &MyResponse{id, "STUFF"}, nil
}
