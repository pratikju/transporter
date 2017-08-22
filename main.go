package main

import "flag"

var (
	mode       = flag.String("mode", "receiver", "Choose the mode of application(sender/receiver)")
	path       = flag.String("path", "", "Absolute path of the file/directory to be transferred")
	recursive  = flag.Bool("R", false, "Whether to send files recursively")
	senderAddr = flag.String("s", "", "Sender Host address to connect to")
	file       = flag.String("pinFileName", "/tmp/password", "name of the file to store network pin.")
)

func main() {
	flag.Parse()
	if *mode == "sender" {
		startServer("0.0.0.0:9999", *path)
	} else {
		startClient(*senderAddr)
	}
}
