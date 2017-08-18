package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

func startClient() {
	senderIP := getSenderIP()
	conn, err := net.Dial("tcp", senderIP.String()+":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	receiveFile(conn)
}

func receiveFile(conn net.Conn) {
	defer conn.Close()

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	for {
		conn.Read(bufferFileSize)
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

		if fileSize == 0 {
			break
		}

		conn.Read(bufferFileName)
		fileName := strings.Trim(string(bufferFileName), ":")

		start := time.Now()

		receivedFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}

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

		receivedFile.Close()

		elapsed := time.Since(start)
		fmt.Println("Transferred : ", fileName)
		fmt.Println("Size        : ", bytefmt.ByteSize(uint64(fileSize)))
		fmt.Println("Time taken  : ", elapsed)
		fmt.Println("-------------------------------------------------------------------------")

	}
}
