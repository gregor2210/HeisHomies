package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Function to receive UDP messages
func udpReceiver(listenPort int, done chan bool) {
	// Create a UDP listener
	addr := fmt.Sprintf(":%d", listenPort)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Printf("Error setting up UDP listener: %v\n", err)
		done <- true
		return
	}
	defer conn.Close()

	fmt.Printf("Listening for UDP responses on port %d...\n", listenPort)

	buffer := make([]byte, 1024)
	for {
		select {
		case <-done:
			fmt.Println("Receiver shutting down...")
			return
		default:
			n, remoteAddr, err := conn.ReadFrom(buffer)
			if err != nil {
				fmt.Printf("Error reading UDP message: %v\n", err)
				continue
			}

			message := strings.TrimSpace(string(buffer[:n]))
			fmt.Printf("Received from %s: %s\n", remoteAddr.String(), message)
		}
	}
}

// Function to send UDP messages
func udpSender(targetIP string, targetPort int, done chan bool) {
	// Create a UDP connection for sending
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", targetIP, targetPort))
	if err != nil {
		fmt.Printf("Error setting up UDP sender: %v\n", err)
		done <- true
		return
	}
	defer conn.Close()

	fmt.Printf("Sending UDP packets to %s:%d...\n", targetIP, targetPort)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a message to send (or 'exit' to quit): ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			fmt.Println("Exiting sender...")
			done <- true
			return
		}

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending UDP message: %v\n", err)
			continue
		}
		fmt.Printf("Sent: %s\n", message)
	}
}

func main() {
	var workspaceNumber int
	fmt.Print("Enter your workspace number: ")
	fmt.Scan(&workspaceNumber)

	// Calculate the port based on the workspace number
	serverPort := 20000 + workspaceNumber
	listenPort := serverPort

	var serverIP string
	fmt.Print("Enter the server IP (or broadcast address): ")
	fmt.Scan(&serverIP)

	// Channel to signal termination
	done := make(chan bool)

	// Start the UDP receiver in a separate goroutine
	go udpReceiver(listenPort, done)

	// Start the UDP sender in a separate goroutine
	go udpSender(serverIP, serverPort, done)

	// Wait for termination signal
	<-done
	close(done)
	fmt.Println("Program exiting...")
}
