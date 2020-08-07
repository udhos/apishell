package main

import (
	"log"
	"net/http"
)

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
