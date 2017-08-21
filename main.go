package main

import "flag"

var (
	mode      = flag.String("mode", "receiver", "Choose the mode of application(sender/receiver)")
	path      = flag.String("path", "", "Absolute path of the file/directory to be transferred")
	recursive = flag.Bool("R", false, "Whether to send files recursively")
)

const (
	port = "9999"
)

func main() {
	flag.Parse()

	if *mode == "sender" {
		startServer("0.0.0.0:"+port, *path)
	} else {
		startClient()
	}
}
