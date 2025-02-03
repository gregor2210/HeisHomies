package main

import (
	"fmt"
	"net"
	"time"
)

const (
	UDP_PORT = ":30000"
	UDP_ADDR = "127.0.0.1"
	TIMEOUT = 3
	TIMEPULSE = 1
)

funpc primary() {
	conn, err := net.Dial("ud", UDP_ADDR + UDP_PORT)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	count := 1

	for {

		fmt.Printf("Primary teller: %d\n", count)
		count++

		_, err = conn.Write([]byte("Hello from primary"))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		time.Sleep(TIMEPULSE * time.Second)


	}

}

func backup() {
	addr, err := net.ResolveUDPAddr("udp", UDP_ADDR + UDP_PORT)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to client")

	buf := make([]byte, 1024)
	lastHeartbeat := time.Now()

	for {
		conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		lastHeartbeat = time.Now()
		fmt.Println("Backup received heartbeat")
	}







func main() {

}

