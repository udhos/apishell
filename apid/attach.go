package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"gopkg.in/yaml.v2"

	"github.com/udhos/apishell/api"
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

	log.Printf("%d %s: url=%s from=%s message: %s", id, me, r.URL.Path, r.RemoteAddr, string(msg))

	var message api.AttachV1Message

	if errYaml := yaml.Unmarshal(msg, &message); errYaml != nil {
		log.Printf("%d %s: url=%s from=%s message: %v", id, me, r.URL.Path, r.RemoteAddr, errYaml)
		return
	}

	log.Printf("%d %s: url=%s from=%s message: %v", id, me, r.URL.Path, r.RemoteAddr, message)

	cmd := exec.Command(message.Args[0], message.Args[1:]...)

	stdin, errStdin := cmd.StdinPipe()
	if errStdin != nil {
		log.Printf("%d %s: url=%s from=%s pipe stdin: %v", id, me, r.URL.Path, r.RemoteAddr, errStdin)
		return
	}
	defer stdin.Close()
	stdout, errStdout := cmd.StdoutPipe()
	if errStdout != nil {
		log.Printf("%d %s: url=%s from=%s pipe stdout: %v", id, me, r.URL.Path, r.RemoteAddr, errStdout)
		return
	}
	defer stdout.Close()
	stderr, errStderr := cmd.StderrPipe()
	if errStderr != nil {
		log.Printf("%d %s: url=%s from=%s pipe stderr: %v", id, me, r.URL.Path, r.RemoteAddr, errStderr)
		return
	}
	defer stderr.Close()

	if errStart := cmd.Start(); errStart != nil {
		log.Printf("%d %s: url=%s from=%s start: %v", id, me, r.URL.Path, r.RemoteAddr, errStart)
		return
	}

	api.WebsocketSpawn(c, io.MultiReader(stdout, stderr), stdin)

	// It is thus incorrect to call Wait before all reads from the pipe have completed.
	// https://golang.org/pkg/os/exec/#Cmd.StdoutPipe
	if errWait := cmd.Wait(); errWait != nil {
		log.Printf("%d %s: url=%s from=%s command wait: %v", id, me, r.URL.Path, r.RemoteAddr, errWait)
	}

	log.Printf("%d %s: url=%s from=%s command exited", id, me, r.URL.Path, r.RemoteAddr)

	/*
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
	*/
}
