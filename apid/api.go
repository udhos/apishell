package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type apiHandler func(w http.ResponseWriter, r *http.Request, app *server)

func registerAPI(app *server, path string, handler apiHandler) {
	log.Printf("registerAPI: registering api: %s", path)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handler %s: url=%s from=%s", path, r.URL.Path, r.RemoteAddr)

		caller := "handler " + path + ":"

		if badBasicAuth(caller, w, r, app) {
			return
		}

		handler(w, r, app)
	})
}

func writeBuf(caller string, w http.ResponseWriter, buf []byte) {
	_, err := w.Write(buf)
	if err != nil {
		log.Printf("%s writeBuf: %v", caller, err)
	}
}

func writeStr(caller string, w http.ResponseWriter, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		log.Printf("%s writeStr: %v", caller, err)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request, app *server) {
	notFound := "404 Nothing here"
	log.Printf("serveRoot: url=%s from=%s %s", r.URL.Path, r.RemoteAddr, notFound)
	http.Error(w, notFound, 404)
}

type v1ExecPayload struct {
	Args []string
}

type response struct {
	HttpStatus int
	ExitStatus int
	Output     string
}

func serveAPIv1Exec(w http.ResponseWriter, r *http.Request, app *server) {
	log.Printf("serveAPIv1Exec: url=%s from=%s", r.URL.Path, r.RemoteAddr)

	sendYaml := false

	badBody := "400 Bad request body"

	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("serveAPIv1Exec: url=%s from=%s body: %v", r.URL.Path, r.RemoteAddr, errRead)
		http.Error(w, badBody, 400)
		return
	}

	var payload v1ExecPayload

	if errYaml := yaml.Unmarshal(body, &payload); errYaml != nil {
		log.Printf("serveAPIv1Exec: url=%s from=%s body: %v", r.URL.Path, r.RemoteAddr, errYaml)
		http.Error(w, badBody, 400)
		return
	}

	if len(payload.Args) < 1 {
		log.Printf("serveAPIv1Exec: url=%s from=%s empty args list", r.URL.Path, r.RemoteAddr)
		http.Error(w, badBody, 400)
		return
	}

	cmd := exec.Command(payload.Args[0], payload.Args[1:]...)
	out, errExec := cmd.CombinedOutput()
	if errExec != nil {

		//serverError := "500 Server error"

		var exitStatus int

		if t, ok := errExec.(*exec.ExitError); ok {
			exitStatus = t.ExitCode()
		}

		log.Printf("serveAPIv1Exec: url=%s from=%s error: %v", r.URL.Path, r.RemoteAddr, errExec)

		//http.Error(w, serverError+": "+errExec.Error(), 500)
		sendResponse(w, 500, exitStatus, out, sendYaml)
		return
	}

	//writeBuf("serveAPIv1Exec", w, out)
	sendResponse(w, 200, 0, out, sendYaml)
}

func sendResponse(w http.ResponseWriter, httpStatus int, exitStatus int, output []byte, sendYaml bool) {
	var result response
	result.HttpStatus = httpStatus
	result.ExitStatus = exitStatus
	result.Output = string(output)

	if sendYaml {
		buf, errMarshal := yaml.Marshal(&result)
		if errMarshal != nil {
			log.Printf("sendResponse yaml.Marshal: %v", errMarshal)
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		httpError(w, string(buf), httpStatus)
		return
	}

	buf, errMarshal := json.Marshal(&result)
	if errMarshal != nil {
		log.Printf("sendResponse json.Marshal: %v", errMarshal)
	}
	w.Header().Set("Content-Type", "application/json")
	httpError(w, string(buf), httpStatus)
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
