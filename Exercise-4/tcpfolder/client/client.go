package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	// Connect to the TCP server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the server. Type messages to send.")

	//reader := bufio.NewReader(os.Stdin)

	// Set TCP_NODELAY to disable Nagle algorithm
	tcpConn := conn.(*net.TCPConn)
	err = tcpConn.SetNoDelay(true)
	if err != nil {
		fmt.Println("Error setting TCP_NODELAY:", err)
		os.Exit(1)
	}
	i := 0
	for {
		//fmt.Print("Enter message: ")
		//message, _ := reader.ReadString('\n')
		message := "Nr: " + strconv.Itoa(i)
		i++
		time.Sleep(100 * time.Millisecond)

		buffer := make([]byte, 1024) // Create a buffer of 1024 bytes
		copy(buffer, message)        // Copy the message into the buffer

		// You can pad the message with a custom byte (e.g., space, null byte, etc.)
		// Here, we use a space byte to pad the remaining space in the buffer
		for j := len(message); j < 1024; j++ {
			buffer[j] = ' ' // Pad with space
		}
		// Send the message to the server
		_, err := conn.Write(buffer)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// Read response from the server
		//buffer := make([]byte, 1024)
		//n, err := conn.Read(buffer)
		//if err != nil {
		//fmt.Println("Error reading response:", err)
		//return
		//}
		//fmt.Printf("Server response: %s\n", string(buffer[:n]))
	}
}
