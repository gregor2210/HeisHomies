package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 8080,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	time.Sleep(20 * time.Second)
	buffer := make([]byte, 1024)
	fmt.Println("Client is running at ")
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Received from server: ", string(buffer[:n]))
		//time.Sleep(25 * time.Millisecond)
	}
}
