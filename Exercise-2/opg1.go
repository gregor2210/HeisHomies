package main

import (
	"fmt"
	"net"
	"os"
)

func listenForServer() string {
	port := ":30000"
	addr, err := net.ResolveUDPAddr("udp", port) // ResolveUDPAddr returns an address of UDP end point
	if err != nil {                              // ResolveUDPAddr(network, address string) (*UDPAddr, error)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr) // ListenUDP listens for incoming UDP packets addressed to the local address laddr on the network net.
	if err != nil {
		fmt.Println("Error listening on UDP port:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)                   // buffer to read the incoming data
	_, remoteaddr, err := conn.ReadFromUDP(buf) // ReadFromUDP reads a UDP packet from c, copying the payload into b. It returns the number of bytes copied into b and the return address that was on the packet.
	if err != nil {
		fmt.Println("Error receiving UDP message:", err)
	}

	fmt.Printf("Server IP: %s\n", remoteaddr.IP.String())
	return remoteaddr.IP.String()
}

func main() {
	serverIP := listenForServer()
	fmt.Println("Server IP: ", serverIP)
}