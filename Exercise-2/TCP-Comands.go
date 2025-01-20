// https://okanexe.medium.com/the-complete-guide-to-tcp-ip-connections-in-golang-1216dae27b5a
package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func send_tcp(conn net.Conn) {

	// Send data to the server
	data := []byte("Hello, Server!")
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read and process data from the server
	// ...
	return
}

func receive_tcp(conn net.Conn) {
	receive_buffer := make([]byte, 1024)

	for {
		// Read data from the server
		n, err := conn.Read(receive_buffer)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Convert to string and find null terminator
		message := string(receive_buffer[:n])
		//for hondtering av null index fra oppgaven

		//------ kan fjærnes basert på om vi bruker fixed eller null fra oppgaven
		nullIndex := strings.Index(message, "\x00")

		if nullIndex != -1 {
			// Trim at first null character
			message = message[:nullIndex]
		}
		//-----

		fmt.Printf("Received: %s\n", message)
	}
}

func main() {
	// Connect to the server
	port := "12345"
	ip := "10.100.23.204" + ":" + port
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	//setter opp nodelay delen fra oppgaven
	// Disable TCP coalescing
	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		tcpConn.SetNoDelay(true) // Disable Nagle's Algorithm
	}

	send_tcp(conn)

	go receive_tcp(conn)

	// Prevent `main` from exiting
	time.Sleep(10 * time.Second) // or use `select{}` to wait forever
}
