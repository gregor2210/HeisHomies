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
			fmt.Println("write failed: ")
			fmt.Println("Error: ", err)
			return
		}
		count++
		time.Sleep(TIMEPULSE * time.Second)


	}

}





func main() {

	primary(0) // Kjører primary i en egen tråd
}

