package main

import (
	"fmt"
	"log"

	"github.com/xtaci/kcp-go/v5"
)

func main() {
	// Listen for incoming connections on the specified address and port
	listenAddr := "localhost:9000"
	conn, err := kcp.Listen(listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is listening on", listenAddr)

	for {
		// Accept new client connections
		client, err := conn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleClient(client.(*kcp.UDPSession))
	}
}

func handleClient(client *kcp.UDPSession) {
	defer client.Close()
	buf := make([]byte, 1024)

	for {
		// Read message from the client
		n, err := client.Read(buf)
		if err != nil {
			log.Println("Error reading from client:", err)
			return
		}
		fmt.Printf("Received message: %s\n", string(buf[:n]))
	}
}
