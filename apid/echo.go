package main

import (
	"log"
	"net/http"
)

func serveAPIEchoV1(w http.ResponseWriter, r *http.Request, app *server, id uint64) {

	me := "serveAPIEchoV1"

	log.Printf("%d %s: url=%s from=%s", id, me, r.URL.Path, r.RemoteAddr)

	c, errUpgrade := upgrader.Upgrade(w, r, nil)
	if errUpgrade != nil {
		log.Printf("%d %s: url=%s from=%s upgrade: %v", id, me, r.URL.Path, r.RemoteAddr, errUpgrade)
		return
	}
	defer c.Close()
	for {
		mt, message, errRead := c.ReadMessage()
		if errRead != nil {
			log.Printf("%d %s: url=%s from=%s read: %v", id, me, r.URL.Path, r.RemoteAddr, errRead)
			break
		}
		log.Printf("%d %s: url=%s from=%s message: %s", id, me, r.URL.Path, r.RemoteAddr, message)
		errWrite := c.WriteMessage(mt, message)
		if errWrite != nil {
			log.Printf("%d %s: url=%s from=%s write: %v", id, me, r.URL.Path, r.RemoteAddr, errWrite)
			break
		}
	}
}
