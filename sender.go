package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"code.cloudfoundry.org/bytefmt"
)

// BUFFERSIZE is the smallest buffer size to be transferred over tcp network
const BUFFERSIZE = 1024

func startServer(addr, path string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	generatePassword()
	fmt.Println("Waiting for the client to authenticate...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		pin := readPin(*file)
		status := isPasswordOK(pin, path, conn)
		if status == false {
			conn.Close()
			fmt.Println("Connection has been dropped due to the incorrect PIN.")
		} else {
			fmt.Println("Welcome to the network!. Your Share is our Care!! Enjoy!!")
			conn.Write([]byte(PASS + ":"))
			if *recursive == true {
				fileList := []string{}
				_ = filepath.Walk(path, func(fpath string, f os.FileInfo, err error) error {
					if fpath == path {
						return nil
					}
					fileList = append(fileList, fpath)
					return nil
				})

				for _, file := range fileList {
					sendFile(conn, file)
				}
				conn.Write(make([]byte, BUFFERSIZE))
				conn.Close()
			}
		}
	}
}

func isPasswordOK(pin, path string, conn net.Conn) bool {
	count := 1
	for ; count <= 3; count++ {
		buffer := make([]byte, 64)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}
		buffer = bytes.Trim(buffer, "\x00")
		if string(buffer[:len(buffer)]) != pin {
			conn.Write([]byte(FAIL + ":"))
			if count == 3 {
				return false
			}
		} else {
			return true
		}
	}
	return false
}

func sendFile(conn net.Conn, filePath string) {
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

	fmt.Println("Transferring : ", fileInfo.Name())
	fmt.Println("Size         : ", bytefmt.ByteSize(uint64(fileInfo.Size())))
	fmt.Println("-------------------------------------------------------------------------")

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

func readPin(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
