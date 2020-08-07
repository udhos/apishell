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
	log.Printf("staticHandler.ServeHTTP url=%s from=%s", r.URL.Path, r.RemoteAddr)

	if badBasicAuth("staticHandler.ServeHTTP", w, r, h.app) {
		return
	}

	h.h.ServeHTTP(w, r)
}

func registerStatic(app *server, path, dir string) {
	log.Printf("mapping www path %s to directory %s", path, dir)
	http.Handle(path, httpHandler{app, http.StripPrefix(path, http.FileServer(http.Dir(dir)))})
}

func badBasicAuth(caller string, w http.ResponseWriter, r *http.Request, app *server) bool {

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	username, password, authOK := r.BasicAuth()
	if !authOK {
		log.Printf("%s url=%s from=%s Basic Auth missing", caller, r.URL.Path, r.RemoteAddr)
		http.Error(w, "401 Unauthenticated", 401)
		return true
	}

	if !app.auth(username, password) {
		log.Printf("%s url=%s from=%s Basic Auth failed", caller, r.URL.Path, r.RemoteAddr)
		http.Error(w, "401 Unauthenticated", 401)
		return true
	}

	log.Printf("%s url=%s from=%s Basic Auth ok", caller, r.URL.Path, r.RemoteAddr)

	return false
}
