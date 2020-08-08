package main

import (
	"log"
	"net/http"
)

type httpHandler struct {
	app *server
	h   http.Handler
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := h.app.next()

	log.Printf("%d staticHandler.ServeHTTP url=%s from=%s", id, r.URL.Path, r.RemoteAddr)

	if badBasicAuth("staticHandler.ServeHTTP", id, w, r, h.app) {
		return
	}

	h.h.ServeHTTP(w, r)
}

func registerStatic(app *server, path, dir string) {
	log.Printf("mapping www path %s to directory %s", path, dir)
	http.Handle(path, httpHandler{app, http.StripPrefix(path, http.FileServer(http.Dir(dir)))})
}
