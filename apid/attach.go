package main

import (
	"log"
	"net/http"

	"github.com/udhos/apishell/api"
	"gopkg.in/yaml.v2"
)

func serveAPIAttachV1(w http.ResponseWriter, r *http.Request, app *server, id uint64) {

	me := "serveAPIAttachV1"

	log.Printf("%d %s: url=%s from=%s", id, me, r.URL.Path, r.RemoteAddr)

	c, errUpgrade := upgrader.Upgrade(w, r, nil)
	if errUpgrade != nil {
		log.Printf("%d %s: url=%s from=%s upgrade: %v", id, me, r.URL.Path, r.RemoteAddr, errUpgrade)
		return
	}
	defer c.Close()

	_, msg, errRead := c.ReadMessage()
	if errRead != nil {
		log.Printf("%d %s: url=%s from=%s read: %v", id, me, r.URL.Path, r.RemoteAddr, errRead)
		return
	}

	var message api.AttachV1Message

	if errYaml := yaml.Unmarshal(msg, &message); errYaml != nil {
		log.Printf("%d %s: url=%s from=%s message: %v", id, me, r.URL.Path, r.RemoteAddr, errYaml)
		return
	}

	log.Printf("%d %s: url=%s from=%s message: %v", id, me, r.URL.Path, r.RemoteAddr, message)

	for {
		mt, m, errRead := c.ReadMessage()
		if errRead != nil {
			log.Println("read:", errRead)
			break
		}
		log.Printf("recv: %s", m)
		errWrite := c.WriteMessage(mt, m)
		if errWrite != nil {
			log.Println("write:", errWrite)
			break
		}
	}
}
