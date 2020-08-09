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
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/udhos/apishell/api"
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

func serveAPIExecV1(w http.ResponseWriter, r *http.Request, app *server, id uint64) {
	me := "serveAPIExecV1"

	log.Printf("%d %s: url=%s from=%s", id, me, r.URL.Path, r.RemoteAddr)

	sendYaml := false

	badBody := "400 Bad request body"

	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("%d %s: url=%s from=%s body: %v", id, me, r.URL.Path, r.RemoteAddr, errRead)
		http.Error(w, badBody, 400)
		return
	}

	var payload api.ExecV1RequestBody

	if errYaml := yaml.Unmarshal(body, &payload); errYaml != nil {
		log.Printf("%d %s: url=%s from=%s body: %v", id, me, r.URL.Path, r.RemoteAddr, errYaml)
		http.Error(w, badBody, 400)
		return
	}

	if len(payload.Args) < 1 {
		log.Printf("%d %s: url=%s from=%s empty args list", id, me, r.URL.Path, r.RemoteAddr)
		http.Error(w, badBody, 400)
		return
	}

	cmd := exec.Command(payload.Args[0], payload.Args[1:]...)

	if payload.Stdin != "" {
		var data []byte

		if strings.HasPrefix(payload.Stdin, api.PrefixBase64) {
			// prefixed with base64: encoded
			s := payload.Stdin[len(api.PrefixBase64):]
			d, errDecode := base64.StdEncoding.DecodeString(s)
			if errDecode != nil {
				log.Printf("%d %s: url=%s from=%s stdin decode: %v", id, me, r.URL.Path, r.RemoteAddr, errDecode)
				http.Error(w, "400 stdin decode", 500)
				return
			}
			data = d
		} else {
			// not encoded
			data = []byte(payload.Stdin)
		}

		stdin, errStdinPipe := cmd.StdinPipe()
		if errStdinPipe != nil {
			log.Printf("%d %s: url=%s from=%s cmd stdin: %v", id, me, r.URL.Path, r.RemoteAddr, errStdinPipe)
			http.Error(w, "500 command input", 500)
			return
		}

		go func() {
			defer stdin.Close()
			n, errWrite := stdin.Write(data)
			if errWrite != nil {
				log.Printf("%d %s: url=%s from=%s cmd stdin write len=%d: %v", id, me, r.URL.Path, r.RemoteAddr, n, errWrite)
			}
		}()
	}

	out, errExec := cmd.CombinedOutput()
	if errExec != nil {
		var exitStatus int

		if t, ok := errExec.(*exec.ExitError); ok {
			exitStatus = t.ExitCode()
		}

		log.Printf("%d %s: url=%s from=%s error: %v", id, me, r.URL.Path, r.RemoteAddr, errExec)

		sendResponse(w, 500, exitStatus, out, errExec.Error(), sendYaml, id)
		return
	}

	sendResponse(w, 200, 0, out, "", sendYaml, id)
}

func sendResponse(w http.ResponseWriter, HTTPStatus int, exitStatus int, output []byte, execError string, sendYaml bool, id uint64) {
	var result api.ExecV1ResponseBody
	result.HTTPStatus = HTTPStatus
	result.ExitStatus = exitStatus
	result.Output = api.PrefixBase64 + base64.StdEncoding.EncodeToString([]byte(output))
	result.Error = execError

	//log.Printf("%d sendResponse: output: %s", id, string(output))

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
