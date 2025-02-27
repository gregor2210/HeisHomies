package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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

		// Set TCP_NODELAY to disable Nagle algorithm
		tcpConn := conn.(*net.TCPConn)
		err = tcpConn.SetNoDelay(true)
		if err != nil {
			fmt.Println("Error setting TCP_NODELAY:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer to read exactly 1024 bytes
	buffer := make([]byte, 1024)
	for {
		// Read exactly 1024 bytes
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// Remove the padding (e.g., spaces) from the message
		message := string(buffer[:n])
		message = strings.TrimSpace(message) // Remove spaces or other padding characters

		// Handle the actual message (e.g., print it)
		fmt.Printf("Received: %s\n", message)

		// Echo the message back to the client
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}
