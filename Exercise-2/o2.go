package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func tcpFixedSizeClient(serverIP string, serverPort int) {
	// Koble til serveren
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server with fixed-size messages")

	// Buffer for meldinger (fast størrelse)
	buffer := make([]byte, 1024)

	// Les velkomstmeldingen
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading welcome message: %v\n", err)
		return
	}
	fmt.Printf("Welcome message: %s\n", string(buffer[:n]))

	// Skriver meldinger til serveren
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message (or 'exit' to quit): ")
		message, _ := reader.ReadString('\n')
		message = message[:len(message)-1] // Fjern linjeskift

		if message == "exit" {
			fmt.Println("Exiting...")
			return
		}

		// Send meldingen til serveren
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			continue
		}

		// Les svaret (fast størrelse)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}
		fmt.Printf("Server response: %s\n", string(buffer[:n]))
	}
}

func main() {
	serverIP := "10.100.23.204"
	serverPort := 34933

	tcpFixedSizeClient(serverIP, serverPort)
}
