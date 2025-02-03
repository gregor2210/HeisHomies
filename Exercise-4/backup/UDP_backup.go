package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	"timeout/timeout"
)

func listen_for_master_and_report(conn *net.UDPConn, nr_receved chan string) {
	buffer := make([]byte, 1024)
	fmt.Println("Client is running at ")

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		num_res := string(buffer[:n])
		nr_receved <- num_res

	}

}

func while_master_is_alive(timerTimeoutChan chan bool, nr_receved chan string) int {
	last_nr := 0
	for {
		select {
		case a := <-timerTimeoutChan:
			if a {
				fmt.Println("Timer timed out")
				timeout.TimerStop()
				return last_nr
			}

		case a := <-nr_receved:
			fmt.Println("Received from master: ", a)
			nr, err := strconv.Atoi(a)
			if err != nil {
				fmt.Println("Conversion error:", err)
				os.Exit(1)
			}
			last_nr = nr
			timeout.TimerStart(3)
		}
	}
}

func main() {
	// Channel to receive timer timeout events
	timerTimeoutChan := make(chan bool)
	nr_receved := make(chan string)

	addr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 8080,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Starting go rutines")
	go timeout.PollTimerTimeout(timerTimeoutChan)
	go listen_for_master_and_report(conn, nr_receved)

	last_nr := while_master_is_alive(timerTimeoutChan, nr_receved)
	fmt.Println("Master died, starting backup from: ", last_nr)

	num_str := "Hei"
	for i := last_nr + 1; i <= 1000; i++ {
		num_str = strconv.Itoa(i)
		time.Sleep(1 * time.Second)
		fmt.Println("Backup sending: ", num_str)
	}

}
