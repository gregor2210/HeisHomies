package main

import (
	"fmt"
	"net"
	"time"
)

const (
	snd_addr = "127.0.0.1:40000"							// IP address of the primary
	rcv_addr = ":40000"								// IP address of the backup			
	TIMEOUT = 3						// Timeout for the backup
	TIMEPULSE = 1					// Time between heartbeats
)

func primary() {
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

	count := 1

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

func backup() {
	rcv_UDP_addr, err := net.ResolveUDPAddr("udp", rcv_addr)            // Lager motta adresse for UDP
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", rcv_UDP_addr)                    // Lager lyttekanal for motta adressen
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to client")

	buffer := make([]byte, 1024)  
	for {
		conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))            // Setter timeout for motta adressen
		n, _, err := conn.ReadFromUDP(buffer)		                          // Leser fra bufferen	
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Backup received number: %s\n", string(buffer[:n]))
	
		time.Sleep(TIMEPULSE * time.Second)
	}

}




func main() {

	go primary() // Kjører primary i en egen tråd
	backup() // Kjører backup i hovedtråden
}

