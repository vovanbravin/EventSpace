package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {

	proxies := map[string]*httputil.ReverseProxy{
		"/v1/auth/": httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http",
			Host: "auth-service:8081"}),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for prefix, proxy := range proxies {
			if strings.HasPrefix(r.URL.Path, prefix) {
				log.Printf("Proxy %s\n", r.URL.Path)
				proxy.ServeHTTP(w, r)
				return
			}
		}
		log.Printf("No found url: %s\n", r.URL.Path)
		http.NotFound(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err.Error())
	}

}
