package main

import (
	"io"
	"log"
	"net/http"
)

type apiHandler func(w http.ResponseWriter, r *http.Request, app *server)

func registerAPI(app *server, path string, handler apiHandler) {
	log.Printf("registerAPI: registering api: %s", path)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handler %s: url=%s from=%s", path, r.URL.Path, r.RemoteAddr)

		caller := "handler " + path + ":"

		if badBasicAuth(caller, w, r, app) {
			return
		}

		handler(w, r, app)
	})
}

func serveRoot(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveRoot: url=%s from=%s", r.URL.Path, r.RemoteAddr)
	http.Error(w, "404 Nothing here", 404)
}

func serveAPI(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveApi: url=%s from=%s", r.URL.Path, r.RemoteAddr)
	writeStr("serveApi", w, "ok\n")
}

func writeBuf(caller string, w http.ResponseWriter, buf []byte) {
	_, err := w.Write(buf)
	if err != nil {
		log.Printf("%s writeBuf: %v", caller, err)
	}
}

func writeStr(caller string, w http.ResponseWriter, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		log.Printf("%s writeStr: %v", caller, err)
	}
}
