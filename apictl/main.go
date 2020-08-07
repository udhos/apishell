package main

import (
	"log"
	"runtime"
)

const (
	me      = "apictl"
	version = "0.0"
)

func main() {
	log.Printf("%s: version=%s runtime=%s GOMAXPROCS=%d ARCH=%s OS=%s", me, version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOARCH, runtime.GOOS)
}
