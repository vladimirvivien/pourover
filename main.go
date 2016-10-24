package main

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "flag"
    "fmt"
    "os"
)

var (
    addr *string
    proxyUrl *string
    proxy *httputil.ReverseProxy
)

func init() {
    addr = flag.String("addr", ":9191", "pourover host address")
    proxyUrl = flag.String("url", "", "remote proxied root URL")
}

func forward(w http.ResponseWriter, r *http.Request){
    proxy.ServeHTTP(w, r)
}

func main() {
    flag.Parse()

    if *proxyUrl == "" {
        fmt.Println("missing url flag, see --help for usage")
        os.Exit(1)
    }

    u, err := url.Parse(*proxyUrl)
    if err !=nil {
        fmt.Println(err)
        os.Exit(1)
    }

    proxy = httputil.NewSingleHostReverseProxy(u)
    mux := http.NewServeMux()
    mux.HandleFunc("/", forward)
    s := &http.Server{
    	Addr:           *addr,
    	Handler:        mux,
    }

    fmt.Printf("pourover proxy addr %v for target url %v\n", *addr, u.String())
    s.ListenAndServe()
}
