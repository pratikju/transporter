package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icrowley/fake"
)

const (
	multicastAddr   = "224.0.0.1:7070" //all-hosts group
	maxDatagramSize = 8192
)

var (
	sigs = make(chan os.Signal, 1)
	done = make(chan bool, 1)
)

func sendMulticastPacket() {
	name := fake.City()
	fmt.Printf("----------------------------------------- %s started -----------------------------------------\n", name)
	fmt.Println("Waiting for receivers....")
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal("ResolveUDPAddr failed: ", err)
	}
	c, _ := net.DialUDP("udp", nil, addr)
	for {
		c.Write([]byte(name))
		time.Sleep(1 * time.Second)
	}
}

func waitForSignal() {
	signal.Notify(sigs, syscall.SIGQUIT)
	<-sigs
	fmt.Println()
	done <- true
}

func scanNetwork() map[string]string {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal("ResolveUDPAddr failed: ", err)
	}

	fmt.Println("Scanning.... Press Ctrl+\\ to select sender")
	udpConn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("ListenMulticastUDP failed: ", err)
	}
	udpConn.SetReadBuffer(maxDatagramSize)

	sendersList := listSenders(udpConn)
	if err := udpConn.Close(); err != nil {
		log.Fatal("UDP connection close failed: ", err)
	}
	return sendersList
}

func listSenders(udpConn *net.UDPConn) map[string]string {
	sendersList := map[string]string{}
	go waitForSignal()
	for {
		select {
		case <-done:
			return sendersList
		default:
			packet := make([]byte, maxDatagramSize)
			_, src, err := udpConn.ReadFromUDP(packet)
			if err != nil {
				log.Fatal("ReadFromUDP failed:", err)
			}
			senderName := string(bytes.Trim(packet, "\x00"))
			if src != nil && sendersList[senderName] == "" {
				fmt.Printf("%s\n", senderName)
				sendersList[senderName] = src.IP.String()
			}
		}
	}
}

func selectSender(sendersList map[string]string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter sender's name: ")
	for scanner.Scan() {
		selectedSender := scanner.Text()
		if _, ok := sendersList[selectedSender]; ok {
			return selectedSender
		}
		fmt.Print("Sender not in list. Enter a sender from the list: ")
	}
	return ""
}

func getSenderIP() string {
	sendersList := scanNetwork()
	senderName := selectSender(sendersList)
	return sendersList[senderName]
}

func padString(source string, toLength int) string {
	currLength := len(source)
	remLength := toLength - currLength

	for i := 0; i < remLength; i++ {
		source += ":"
	}
	return source
}
