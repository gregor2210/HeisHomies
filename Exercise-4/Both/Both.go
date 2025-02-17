package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func send_udp_message(conn *net.UDPConn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

func master(lastCount int, addr *net.UDPAddr) {
	//cmd := exec.Command("cmd", "/C", "start", "cmd", "/K", "go run Both.go")
	cmd := exec.Command("gnome-terminal", "--", "go", "run", "Both.go")

	// Start the new process
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting new terminal:", err)
	}
	fmt.Println("New terminal started")

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	num := -1 //deafult number
	//buffer :=make([]byte, 1024)
	num_str := strconv.Itoa(num)

	for i := lastCount; i <= 1000; i++ {
		num_str = strconv.Itoa(i)
		send_udp_message(conn, num_str)
		fmt.Println("Sent from Master: ", num_str)
		time.Sleep(1 * time.Second)
	}
}

//------------------------------------------------------------

func backup(addr *net.UDPAddr) int {
	lastCount := 0
	TIMEOUT := 3
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	buffer := make([]byte, 1024)

	//Listening and setting timer
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT) * time.Second)) // Setter timeout for motta adressen
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Primary process not responding, taking over as primary")
			break
		}
		num_res := string(buffer[:n])

		nr, err := strconv.Atoi(num_res)
		if err != nil {
			fmt.Println("Conversion error:", err)
			os.Exit(1)
		}
		lastCount = nr
		fmt.Println("Received from master: ", num_res)

	}

	return lastCount

}

func main() {

	addr := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 8080,
	}

	lastCount := 0
	lastCount = backup(addr)
	fmt.Println("Master died, starting backup from: ", lastCount)
	master(lastCount, addr)

}
