package main

import (
	"log"
	"net"
	"time"
)

const (
	multicastAddr   = "224.0.0.1:7070" //all-hosts group
	maxDatagramSize = 8192
)

func sendMulticastPacket() {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal("ResolveUDPAddr failed: ", err)
	}
	c, _ := net.DialUDP("udp", nil, addr)
	for {
		c.Write([]byte(""))
		time.Sleep(1 * time.Second)
	}
}

func getSenderIP() net.IP {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal("ResolveUDPAddr failed: ", err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("ListenMulticastUDP failed: ", err)
	}
	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		_, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		if err := l.Close(); err != nil {
			log.Fatal("UDP connection close failed: ", err)
		}
		return src.IP
	}
}

func padString(source string, toLength int) string {
	currLength := len(source)
	remLength := toLength - currLength

	for i := 0; i < remLength; i++ {
		source += ":"
	}
	return source
}
