package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func receiveFile(conn net.Conn) {
	defer conn.Close()

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	log.Println(fileName)

	receivedFile, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
	}
	defer receivedFile.Close()

	var receivedBytes int64
	var remainingSize int64
	for {
		remainingSize = fileSize - receivedBytes
		if remainingSize < BUFFERSIZE {
			io.CopyN(receivedFile, conn, remainingSize)
			conn.Read(make([]byte, (BUFFERSIZE - remainingSize)))
			break
		}

		io.CopyN(receivedFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE

	}

}

func startClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	receiveFile(conn)
}
