package api

import (
	"io"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebsocketSpawn attaches websocket to reader and writer.
func WebsocketSpawn(c *websocket.Conn, r io.Reader, w io.WriteCloser) {
	var wg sync.WaitGroup
	wg.Add(2)
	go copyFromReader(&wg, c, r)
	go copyToWriter(&wg, c, w)
	wg.Wait()
}

func copyFromReader(wg *sync.WaitGroup, c *websocket.Conn, r io.Reader) {
	defer wg.Done()
	defer c.Close()
	buf := make([]byte, 1000)
	//log.Printf("copyFromReader: reading reader into websocket")
	for {
		n, errRead := r.Read(buf)
		//log.Printf("copyFromReader: read: %d", n)
		if errRead != nil {
			log.Printf("copyFromReader: read: %v", errRead)
			if n < 1 {
				return // empty read
			}
		}
		errWrite := c.WriteMessage(websocket.BinaryMessage, buf[:n])
		//log.Printf("copyFromReader: write: %d", n)
		if errWrite != nil {
			log.Printf("copyFromReader: write: %v", errWrite)
			return
		}
	}
}

func copyToWriter(wg *sync.WaitGroup, c *websocket.Conn, w io.WriteCloser) {
	defer wg.Done()
	defer w.Close()
	//log.Printf("copyToWriter: reading from websocket into writer")
	for {
		_, m, errRead := c.ReadMessage()
		n := len(m)
		//log.Printf("copyToWriter: read: %d", n)
		if errRead != nil {
			log.Printf("copyToWriter: read: %v", errRead)
			if n < 1 {
				return // empty read
			}
		}
		_, errWrite := w.Write(m)
		//log.Printf("copyToWriter: write: %d", n2)
		if errWrite != nil {
			log.Printf("copyToWriter: write: %v", errWrite)
			return
		}
	}
}
