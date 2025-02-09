package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	snd_addr  = "127.0.0.1:40000" // IP address of the primary
	rcv_addr  = ":40000"          // IP address of the backup
	TIMEOUT   = 3                 // Timeout for the backup
	TIMEPULSE = 1                 // Time between heartbeats
)

func backup(countChan chan int) {
	rcv_UDP_addr, err := net.ResolveUDPAddr("udp", rcv_addr) // Lager motta adresse for UDP
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", rcv_UDP_addr) // Lager lyttekanal for motta adressen
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to client")

	buffer := make([]byte, 1024)
	lastCount := 0

	for {
		conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second)) // Setter timeout for motta adressen
		n, _, err := conn.ReadFromUDP(buffer)                       // Leser fra bufferen
		if err != nil {
			countChan <- lastCount
			fmt.Println("Primary process not responding, taking over as primary")
			return
		}
		fmt.Printf("Backup received number: %s\n", string(buffer[:n]))

		received := string(buffer[:n])
		count, err := strconv.Atoi(received)
		if err != nil {
			fmt.Println("Conversion error:", err)
			continue
		}
		lastCount = count
		fmt.Printf("Backup received number: %d\n", count)

	}

}

func main() {

	countChan := make(chan int) // Kjører backup i hovedtråden

	go backup(countChan)

	lastCount := <-countChan
	fmt.Printf("Backup received number: %d\n", lastCount)

	for i := lastCount + 1; i < 100; i++ {
		fmt.Printf("Backup counting: %d\n", i)
		time.Sleep(TIMEPULSE * time.Second)
	}
}
