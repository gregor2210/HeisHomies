package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)
num := -1 //deafult number


func send_udp_message(conn *net.UDPConn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

func timer_addtime(duration float64) {
	//add time to timer
}

func isTimerExpired() bool {
	//retunerer true eller false hvis timeren er expired
}

func listen_for_master_and_report(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	fmt.Println("Client is running at ")

	
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		num_res := string(buffer[:n])
		fmt.Println("Received from master: ", num_res)
		//time.Sleep(25 * time.Millisecond)

		timer_addtime(3)
	}	
}

func main() {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 8080,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	

}
