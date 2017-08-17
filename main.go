package main

import (
	"flag"
	"fmt"
)

var (
	mode      = flag.String("mode", "receiver", "Choose the mode of application(sender/receiver)")
	path      = flag.String("path", "", "Absolute path of the file/directory to be transferred")
	recursive = flag.Bool("R", false, "Whether to send files recursively")
	port      = flag.Int("p", 7080, "Port to connect to")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", *port)

	if *mode == "sender" {
		startServer(addr, *path)
	} else {
		startClient(addr)
	}
}
