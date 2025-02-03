package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func send_udp_message(conn *net.UDPConn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 8080,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	num := -1 //deafult number
	//buffer :=make([]byte, 1024)
	num_str := strconv.Itoa(num)

	for i := 0; i <= 1000; i++ {
		num_str = strconv.Itoa(i)
		send_udp_message(conn, num_str)
		time.Sleep(1 * time.Second)
		fmt.Println("Sent from Master: ", num_str)
	}

}
