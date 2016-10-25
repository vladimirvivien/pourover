package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	addr     *string
	proxyUrl *string
	loglevel *string
	proxy    *httputil.ReverseProxy
)

func init() {
	addr = flag.String("addr", ":9191", "pourover host address")
	proxyUrl = flag.String("url", "", "remote proxied root URL")
	loglevel = flag.String("log", "info", "log level (info, debug supported)")
}

func forward(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s - %s", r.Method, r.Proto, r.RequestURI)
	switch *loglevel {
	case "debug":
		for k, v := range r.Header {
			log.Printf("\t>%s:%v", k, v)
		}
	}
	proxy.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	if *proxyUrl == "" {
		fmt.Println("missing url flag, see --help for usage")
		os.Exit(1)
	}

	u, err := url.Parse(*proxyUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proxy = httputil.NewSingleHostReverseProxy(u)
	mux := http.NewServeMux()
	mux.HandleFunc("/", forward)
	s := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	fmt.Printf("pourover proxy addr %v for target url %v\n", *addr, u.String())
	s.ListenAndServe()
}
