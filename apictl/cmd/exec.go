package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/udhos/apishell/api"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command on apid, the api shell server.",
	Long: `Execute a command on apid, the api shell server.

apict exec [--stdin string|@file] cmd arg1..argN

Example:

apictl exec --stdin "hello world" wc

apictl exec --stdin @/etc/passwd head

apictl exec head /etc/passwd

apictl exec -- ls -la

apictl exec -- bash -c "echo -n 12345 | wc"
`,
	Run: func(cmd *cobra.Command, args []string) {

		var body api.ExecV1RequestBody
		body.Args = args

		if execStdin != "" {
			if execStdin[0] == '@' {
				data, errRead := ioutil.ReadFile(execStdin[1:])
				if errRead != nil {
					fmt.Printf("stdin: %s: %v\n", execStdin, errRead)
					return
				}
				body.Stdin = api.PrefixBase64 + base64.StdEncoding.EncodeToString(data)
			} else {
				body.Stdin = api.PrefixBase64 + base64.StdEncoding.EncodeToString([]byte(execStdin))
			}
		}

		var buf bytes.Buffer

		encoder := json.NewEncoder(&buf)
		if errEnc := encoder.Encode(&body); errEnc != nil {
			fmt.Printf("encode json: %v\n", errEnc)
			return
		}

		//fmt.Printf("request body: %s\n", buf.String())

		user := "admin"
		pass := "admin"
		url := "https://localhost:8080" + api.ExecV1Path
		req, errReq := http.NewRequest("POST", url, &buf)
		if errReq != nil {
			fmt.Printf("request: %s: %v\n", url, errReq)
			return
		}
		req.SetBasicAuth(user, pass)
		c := newHTTPClient(true)
		resp, errPost := c.Do(req)
		if errPost != nil {
			fmt.Printf("post: %s: %v\n", url, errPost)
			return
		}
		if resp.StatusCode != 200 {
			fmt.Printf("StatusCode: %d\n", resp.StatusCode)
			fmt.Printf("Status: %s\n", resp.Status)
			return
		}

		respBody, errBody := ioutil.ReadAll(resp.Body)
		if errBody != nil {
			fmt.Printf("body: %v\n", errBody)
			return
		}

		fmt.Fprintf(os.Stderr, "Body Length: %d\n", len(respBody))

		var result api.ExecV1ResponseBody

		if errUnmarshal := yaml.Unmarshal(respBody, &result); errUnmarshal != nil {
			fmt.Printf("unmarshaml yaml: %v\n", errUnmarshal)
			return
		}

		//fmt.Printf("result: %v\n", result)

		fmt.Fprintf(os.Stderr, "HTTPStatus: %d\n", result.HTTPStatus)
		fmt.Fprintf(os.Stderr, "ExitStatus: %d\n", result.ExitStatus)
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)

		var output string

		if strings.HasPrefix(result.Output, api.PrefixBase64) {
			// prefixed with base64: encoded
			suffix := result.Output[len(api.PrefixBase64):]
			o, errDecode := base64.StdEncoding.DecodeString(suffix)
			if errDecode != nil {
				fmt.Printf("decode base64: %v", errDecode)
				return
			}
			output = string(o)
		} else {
			// not encoded
			output = result.Output
		}

		fmt.Print(output)
	},
}

func newHTTPClient(tlsInsecureSkipVerify bool) http.Client {
	log.Printf("newHTTPClient: tlsInsecureSkipVerify=%v", tlsInsecureSkipVerify)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: tlsInsecureSkipVerify,
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	c := http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
	return c
}

var execStdin string

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	execCmd.Flags().StringVar(&execStdin, "stdin", "", "Input for command. --stdin string OR --stdin @file")
}
