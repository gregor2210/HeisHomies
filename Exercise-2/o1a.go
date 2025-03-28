package main

import (
	"fmt"
	"net"
	"os"
)



func listenForServer() string {

	//LIsten for server via UDP
	port := ":30000"
	addr, err := net.ResolveUDPAddr("udp", port) 
	//error handling:
	if err != nil {                              
		os.Exit(1)
	}
	//Listen for message on port
	conn, err := net.ListenUDP("udp", addr)

	//error handling:
	if err != nil {
		fmt.Println("Error listening on UDP port:", err)
		os.Exit(1)
	}
	defer conn.Close()
	//allocates data for incomming messages and recieving/processing data->message
	buf := make([]byte, 1024)                  
	_, remoteaddr, err := conn.ReadFromUDP(buf)


	//error handling
	if err != nil {
		fmt.Println("Error receiving UDP message:", err)
	}
	//print
	fmt.Printf("Server IP: %s\n", remoteaddr.IP.String())
	return remoteaddr.IP.String()
}

func main() {
	serverIP := listenForServer()
	fmt.Println("Server IP: ", serverIP)
}