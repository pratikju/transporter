package main

import (
	"flag"
	"fmt"
)

var (
	sender    = flag.Bool("s", false, "Specify whether application should act as sender ")
	path      = flag.String("path", "", "Absolute path of the file/directory to be transferred")
	recursive = flag.Bool("R", false, "Whether to send files recursively")
)

const (
	port = 9999
)

func main() {
	flag.Parse()

	if *sender {
		startServer(fmt.Sprintf("0.0.0.0:%d", port), *path)
	} else {
		startClient()
	}
}
