package cmd

import (
	"crypto/tls"
	"encoding/base64"
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

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo [remote command]",
	Short: "Execute echo on apid, the shell server.",
	Long: `Execute echo on apid, the api shell server.

apictl echo

Example:

apictl echo
`,
	Run: func(cmd *cobra.Command, args []string) {

		errlog := log.New(os.Stderr, "", 0)

		host := viper.GetString("server")
		errlog.Printf("apid server host: %s", host)

		u := url.URL{Scheme: "wss", Host: host, Path: api.EchoV1Path}
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
	rootCmd.AddCommand(echoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// echoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// echoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
