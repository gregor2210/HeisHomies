package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

func main() {
	// Connect to the server
	serverAddr := "localhost:9000"
	conn, err := kcp.Dial(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Client connected to", serverAddr)
	i := 0
	for {
		// Send a message to the server every second
		message := "Hello from client!" + strconv.Itoa(i)
		i++
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		fmt.Printf("Sent message: %s\n", message)
		time.Sleep(100 * time.Millisecond)
	}
}
