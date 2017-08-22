package main

import (
	"bufio"
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

const PASS = "PASS"
const FAIL = "FAIL"

func startClient() {
	var response string
	var i int
	senderIP := getSenderIP()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", senderIP.String(), port))
	if err != nil {
		log.Fatalln(err)
	}
	for i = 0; i < 3; i++ {
		fmt.Println("Enter the PIN: ")
		_, err = fmt.Scanln(&response)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(conn, response)
		resp, _ := bufio.NewReader(conn).ReadString(':')
		resp = strings.TrimSuffix(resp, ":")
		if resp != FAIL {
			break
		}
		if i < 2 {
			fmt.Printf("Incorrect Attempt(s).%d attempts more...\n\n", 2-i)
		}

	}

	if i == 3 {
		fmt.Println("Maximum number of incorrect PIN attempts crossed.")
		return
	}

	receiveFile(conn)
}

func receiveFile(conn net.Conn) {
	defer conn.Close()
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	for {
		_, err := conn.Read(bufferFileSize)
		if err == io.EOF {
			fmt.Println("The connection has been dropped..")
			conn.Close()
		}

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
		fmt.Println("-----------------------Receiving File(s)---------------------------------")
		fmt.Println("Transferred : ", fileName)
		fmt.Println("Size        : ", bytefmt.ByteSize(uint64(fileSize)))
		fmt.Println("Time taken  : ", elapsed)
		fmt.Println("-------------------------------------------------------------------------")

	}
}
