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



func primary(startValue int) {
	send_UDP_addr, err := net.ResolveUDPAddr("udp", snd_addr)    // Lager sendeadrese for UDP
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, send_UDP_addr)       // Lager en UDP-tilkobling / server
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()                                   // Lukker tilkoblingen når ferdig og ikke før

	fmt.Println("Connected to server")

	count := startValue

	for {

		fmt.Printf("Primary teller: %d\n", count)


		_, err = conn.Write([]byte(fmt.Sprintf("%d", count)))    // Sender melding til server
		
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		count++
		time.Sleep(TIMEPULSE * time.Second)


	}

}



func main() {
	go primary(0)
	countChan := make(chan int) // Kjører backup i hovedtråden

	go backup(countChan)

	lastCount := <-countChan
	fmt.Printf("Backup received number: %d\n", lastCount)

	for i := lastCount + 1; i < 100; i++ {
		fmt.Printf("Backup counting: %d\n", i)
		time.Sleep(TIMEPULSE * time.Second)
	}
}
