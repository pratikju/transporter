package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

// BUFFERSIZE is size of buffer to transferred at once.
const BUFFERSIZE = 1024

var (
	mode = flag.String("mode", "", "Choose the mode of application")
	path = flag.String("path", "", "Absolute path of the file/directory to be transferred")
	port = flag.Int("p", 7080, "Port to connect to")
)

func padString(source string, toLength int) string {
	currLength := len(source)
	remLength := toLength - currLength

	for i := 0; i < remLength; i++ {
		source += ":"
	}
	return source
}

func sendFile(conn net.Conn, filePath string) {
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}

	fileName := padString(fileInfo.Name(), 64)
	fileSize := padString(strconv.FormatInt(fileInfo.Size(), 10), 10)

	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))

	sendBuffer := make([]byte, BUFFERSIZE)

	for {

		_, err := file.Read(sendBuffer)
		if err != nil {
			if err != io.EOF {
				return
			}
			break
		}

		conn.Write(sendBuffer)
	}
	return

}

func startServer(addr, path string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go sendFile(conn, path)
	}

}

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", *port)

	if *mode == "server" {
		startServer(addr, *path)
	} else {
		startClient(addr)
	}
}
