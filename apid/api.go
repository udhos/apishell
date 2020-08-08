package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"gopkg.in/yaml.v2"
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

type v1ExecPayload struct {
	Stdin string
	Args  []string
}

type response struct {
	HTTPStatus int
	ExitStatus int
	Output     string
	Error      string
}

func serveAPIv1Exec(w http.ResponseWriter, r *http.Request, app *server, id uint64) {
	log.Printf("%d serveAPIv1Exec: url=%s from=%s", id, r.URL.Path, r.RemoteAddr)

	sendYaml := false

	badBody := "400 Bad request body"

	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("%d serveAPIv1Exec: url=%s from=%s body: %v", id, r.URL.Path, r.RemoteAddr, errRead)
		http.Error(w, badBody, 400)
		return
	}

	var payload v1ExecPayload

	if errYaml := yaml.Unmarshal(body, &payload); errYaml != nil {
		log.Printf("%d serveAPIv1Exec: url=%s from=%s body: %v", id, r.URL.Path, r.RemoteAddr, errYaml)
		http.Error(w, badBody, 400)
		return
	}

	if len(payload.Args) < 1 {
		log.Printf("%d serveAPIv1Exec: url=%s from=%s empty args list", id, r.URL.Path, r.RemoteAddr)
		http.Error(w, badBody, 400)
		return
	}

	cmd := exec.Command(payload.Args[0], payload.Args[1:]...)

	if len(payload.Stdin) > 0 {
		data, errDecode := base64.StdEncoding.DecodeString(payload.Stdin)
		if errDecode != nil {
			log.Printf("%d serveAPIv1Exec: url=%s from=%s stdin decode: %v", id, r.URL.Path, r.RemoteAddr, errDecode)
			http.Error(w, "400 stdin decode", 500)
			return
		}

		stdin, errStdinPipe := cmd.StdinPipe()
		if errStdinPipe != nil {
			log.Printf("%d serveAPIv1Exec: url=%s from=%s cmd stdin: %v", id, r.URL.Path, r.RemoteAddr, errStdinPipe)
			http.Error(w, "500 command input", 500)
			return
		}

		go func() {
			defer stdin.Close()
			n, errWrite := stdin.Write(data)
			if errWrite != nil {
				log.Printf("%d serveAPIv1Exec: url=%s from=%s cmd stdin write len=%d: %v", id, r.URL.Path, r.RemoteAddr, n, errWrite)
			}
		}()
	}

	out, errExec := cmd.CombinedOutput()
	if errExec != nil {
		var exitStatus int

		if t, ok := errExec.(*exec.ExitError); ok {
			exitStatus = t.ExitCode()
		}

		log.Printf("%d serveAPIv1Exec: url=%s from=%s error: %v", id, r.URL.Path, r.RemoteAddr, errExec)

		sendResponse(w, 500, exitStatus, out, errExec.Error(), sendYaml, id)
		return
	}

	sendResponse(w, 200, 0, out, "", sendYaml, id)
}

func sendResponse(w http.ResponseWriter, HTTPStatus int, exitStatus int, output []byte, execError string, sendYaml bool, id uint64) {
	var result response
	result.HTTPStatus = HTTPStatus
	result.ExitStatus = exitStatus
	result.Output = string(output)
	result.Error = execError

	if sendYaml {
		buf, errMarshal := yaml.Marshal(&result)
		if errMarshal != nil {
			log.Printf("%d sendResponse yaml.Marshal: %v", id, errMarshal)
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		httpError(w, string(buf), HTTPStatus)
		return
	}

	buf, errMarshal := json.Marshal(&result)
	if errMarshal != nil {
		log.Printf("%d sendResponse json.Marshal: %v", id, errMarshal)
	}
	w.Header().Set("Content-Type", "application/json")
	httpError(w, string(buf), HTTPStatus)
}

// httpError does not reset Content-Type, http.Error does.
// httpError replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
func httpError(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
