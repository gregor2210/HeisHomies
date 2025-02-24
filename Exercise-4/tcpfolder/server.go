package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Start listening on a TCP port
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed")
			return
		}
		fmt.Printf("Received: %s\n", string(buffer[:n]))

		// Echo back the message
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}
