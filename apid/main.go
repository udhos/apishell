package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
)

const (
	me      = "apid"
	version = "0.0"
)

type server struct {
	basicAuthUser string
	basicAuthPass string
}

func main() {
	log.Printf("%s: version=%s runtime=%s GOMAXPROCS=%d ARCH=%s OS=%s", me, version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOARCH, runtime.GOOS)

	var app server

	var key, cert, listen string

	flag.StringVar(&key, "key", "key.pem", "TLS key file")
	flag.StringVar(&cert, "cert", "cert.pem", "TLS cert file")
	flag.StringVar(&listen, "listen", ":8080", "listen address")
	flag.Parse()

	registerAPI(&app, "/api/", serveAPI)

	registerStatic(&app, "/static/", ".")

	log.Printf("serving HTTPS on TCP %s", listen)
	if err := http.ListenAndServeTLS(listen, cert, key, nil); err != nil {
		log.Fatalf("ListenAndServeTLS: %s: %v", listen, err)
	}
}
