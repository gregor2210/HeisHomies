package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	UDP_PORT  = ":5005"
	UDP_ADDR  = "127.0.0.1"
	TIMEOUT   = 3 * time.Second // Backup antar at primary er død etter 3 sekunder
	INTERVAL  = 1 * time.Second // Primary sender heartbeat hvert sekund
)

// Primary-prosessen: Teller og sender heartbeats
func primary() {
	fmt.Println("Primary starter...")

	// Starter en backup-prosess
	spawnBackup()

	conn, err := net.Dial("udp", UDP_ADDR+UDP_PORT)
	if err != nil {
		fmt.Println("Feil ved opprettelse av UDP-socket:", err)
		return
	}
	defer conn.Close()

	count := 1
	for {
		fmt.Printf("Primary teller: %d\n", count)
		count++

		// Send heartbeat til backup
		_, err := conn.Write([]byte("alive"))
		if err != nil {
			fmt.Println("Feil ved sending av heartbeat:", err)
			return
		}
		time.Sleep(INTERVAL)
	}
}

// Backup-prosessen: Lytter etter heartbeats og tar over hvis de stopper
func backup() {
	fmt.Println("Backup venter på heartbeat...")

	addr, err := net.ResolveUDPAddr("udp", UDP_PORT)
	if err != nil {
		fmt.Println("Feil ved ResolveUDPAddr:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Feil ved opprettelse av UDP-server:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 10)
	lastHeartbeat := time.Now()

	for {
		conn.SetReadDeadline(time.Now().Add(TIMEOUT))
		_, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Backup: Ingen heartbeat mottatt, antar at primary er død!")
			primary() // Tar over som ny primary
			return
		}

		fmt.Println("Backup: Mottok heartbeat")
		lastHeartbeat = time.Now()
	}
}

// Starter en ny backup-prosess
func spawnBackup() {
	cmd := exec.Command(os.Args[0], "backup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Feil ved oppstart av backup:", err)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "backup" {
		backup()
	} else {
		primary()
	}
}
