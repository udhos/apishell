package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"sync"
)

const (
	me      = "apid"
	version = "0.0"
)

type server struct {
	basicAuthUser string
	basicAuthPass string
	lock          sync.RWMutex
}

func (s *server) auth(user, pass string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return user == s.basicAuthUser && pass == s.basicAuthPass
}

func main() {
	log.Printf("%s: version=%s runtime=%s GOMAXPROCS=%d ARCH=%s OS=%s", me, version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOARCH, runtime.GOOS)

	var app server
	var key, cert, listen string
	var staticPath, staticDir string

	flag.StringVar(&key, "key", "key.pem", "TLS key file")
	flag.StringVar(&cert, "cert", "cert.pem", "TLS cert file")
	flag.StringVar(&listen, "listen", ":8080", "listen address")
	flag.StringVar(&app.basicAuthUser, "basicAuthUser", "admin", "basic auth username")
	flag.StringVar(&app.basicAuthPass, "basicAuthPass", "admin", "basic auth password")
	flag.StringVar(&staticPath, "staticPath", "/static/", "static path")
	flag.StringVar(&staticDir, "staticDir", ".", "static dir")
	flag.Parse()

	registerStatic(&app, staticPath, staticDir)

	registerAPI(&app, "/", serveRoot)
	registerAPI(&app, "/api/v1/exec/", serveAPIv1Exec)

	if err := listenAndServeTLS(listen, cert, key, nil); err != nil {
		log.Fatalf("ListenAndServeTLS: %s: %v", listen, err)
	}
}

func listenAndServeTLS(listen, certFile, keyFile string, handler http.Handler) error {
	server := &http.Server{Addr: listen, Handler: handler}
	log.Printf("listenAndServeTLS: serving HTTPS on TCP %s", listen)
	return server.ListenAndServeTLS(certFile, keyFile)
}
