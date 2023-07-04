package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

func main() {
	target := "https://httpbin.org"
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	r := mux.NewRouter().Host("{subdomain:[a-z]+}.puzzle.com").Subrouter()

	r.HandleFunc("/target/{forward:.*}", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "https"
		r.URL.Path = mux.Vars(r)["forward"]
		r.Host = mux.Vars(r)["subdomain"] + ".org"

		proxy.ServeHTTP(w, r)
	})
	http.Handle("/", r)

	http.ListenAndServe(":3222", r)
}
