package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func generatePassword() {
	password := String(6)
	fmt.Println("Please distribute this pin to your clients: \n", password)
	data := []byte(password)
	f, err := os.Create(*file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}
