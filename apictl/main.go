package main

import (
	"log"
	"runtime"

	"github/udhos/apishell/apictl/cmd"
)

const (
	me      = "apictl"
	version = "0.0"
)

func main() {
	log.Printf("%s: version=%s runtime=%s GOMAXPROCS=%d ARCH=%s OS=%s", me, version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOARCH, runtime.GOOS)
	cmd.Execute()
}
