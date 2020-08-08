package main

import (
	"log"
	"net/http"
)

func badBasicAuth(caller string, id uint64, w http.ResponseWriter, r *http.Request, app *server) bool {

	username, password, authOK := r.BasicAuth()
	if !authOK {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		log.Printf("%d %s url=%s from=%s Basic Auth missing", id, caller, r.URL.Path, r.RemoteAddr)
		http.Error(w, "401 Unauthenticated", 401)
		return true
	}

	log.Printf("%d %s url=%s from=%s Basic Auth username=%s", id, caller, r.URL.Path, r.RemoteAddr, username)

	if !app.auth(username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		log.Printf("%d %s url=%s from=%s Basic Auth failed", id, caller, r.URL.Path, r.RemoteAddr)
		http.Error(w, "401 Unauthenticated", 401)
		return true
	}

	log.Printf("%d %s url=%s from=%s Basic Auth ok", id, caller, r.URL.Path, r.RemoteAddr)

	return false
}
