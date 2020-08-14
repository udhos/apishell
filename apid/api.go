package main

import (
	"io"
	"log"
	"net/http"
)

type apiHandler func(w http.ResponseWriter, r *http.Request, app *server, id uint64)

func registerAPI(app *server, path string, handler apiHandler) {
	log.Printf("registerAPI: registering api: %s", path)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		id := app.next()

		log.Printf("%d handler %s: url=%s from=%s", id, path, r.URL.Path, r.RemoteAddr)

		caller := "handler " + path + ":"

		if badBasicAuth(caller, id, w, r, app) {
			return
		}

		handler(w, r, app, id)
	})
}

func writeBuf(caller string, id uint64, w http.ResponseWriter, buf []byte) {
	_, err := w.Write(buf)
	if err != nil {
		log.Printf("%d %s writeBuf: %v", id, caller, err)
	}
}

func writeStr(caller string, id uint64, w http.ResponseWriter, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		log.Printf("%d %s writeStr: %v", id, caller, err)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request, app *server, id uint64) {
	notFound := "404 Nothing here"
	log.Printf("%d serveRoot: url=%s from=%s %s", id, r.URL.Path, r.RemoteAddr, notFound)
	http.Error(w, notFound, 404)
}
