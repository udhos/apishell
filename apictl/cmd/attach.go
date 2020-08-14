package cmd

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/udhos/apishell/api"
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach [remote command]",
	Short: "Execute a command on apid, the api shell server.",
	Long: `Execute a command on apid, the api shell server.

apictl attach [--stdin string|@file] cmd arg1..argN

Example:

apictl attach cat
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("attach called")

		errlog := log.New(os.Stderr, "", 0)

		host := viper.GetString("server")
		errlog.Printf("apid server host: %s", host)

		u := url.URL{Scheme: "wss", Host: host, Path: api.AttachV1Path}
		log.Printf("connecting to %s", u.String())

		user := "admin"
		pass := "admin"
		h := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))}}

		tlsInsecureSkipVerify := true
		errlog.Printf("tlsInsecureSkipVerify=%v", tlsInsecureSkipVerify)

		dialer := &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: tlsInsecureSkipVerify,
			},
		}
		c, _, err := dialer.Dial(u.String(), h)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		var message api.AttachV1Message
		message.Args = args

		buf, errMarshal := json.Marshal(&message)
		if errMarshal != nil {
			errlog.Printf("json marshal: %v", errMarshal)
			return
		}

		errWrite := c.WriteMessage(websocket.TextMessage, buf)
		if errWrite != nil {
			errlog.Printf("websocket write: %v", errWrite)
			return
		}

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: mt=%d %s", mt, message)
			}
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// attachCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// attachCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
