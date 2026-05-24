package main

import (
	"github.com/rs/cors"
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

	fs := http.FileServer(http.Dir("./web"))

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		for prefix, proxy := range proxies {
			if strings.HasPrefix(r.URL.Path, prefix) {
				log.Printf("Proxy %s %s", r.Method, r.URL.Path)
				proxy.ServeHTTP(w, r)
				return
			}
		}
		fs.ServeHTTP(w, r)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", mainHandler)

	handler := cors.Default().Handler(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err.Error())
	}

}
