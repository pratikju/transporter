package main

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

// BUFFERSIZE is the smallest buffer size to be transferred over tcp network
const BUFFERSIZE = 1024

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

		if *recursive == true {

			fileList := []string{}
			_ = filepath.Walk(path, func(fpath string, f os.FileInfo, err error) error {
				if fpath == path {
					return nil
				}
				fileList = append(fileList, fpath)
				return nil
			})

			log.Println(fileList)

			for _, file := range fileList {
				sendFile(conn, file)
			}
			conn.Write(make([]byte, BUFFERSIZE))
			conn.Close()
		}
	}

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
