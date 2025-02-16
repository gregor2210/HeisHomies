package connectivity

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	// Timeout for receiving UDP messages
	TIMEOUT = 3
	// World view sending UDP connection setup
	UDP_world_view_send_port = 8080
	UDP_world_view_send_ip   = "127.0.0.1"

	// World view receiving UDP connection setup
	UDP_world_view_receive_port = 8070
	UDP_world_view_receive_ip   = "127.0.0.1"
)

var (
	// World view sending UDP connection setup
	addr_sending_world_view *net.UDPAddr
	conn_sending_world_view *net.UDPConn

	// World view receiving UDP connection setup
	addr_receiving_world_view *net.UDPAddr
	conn_receiving_world_view *net.UDPConn
)

// // World view sending UDP connection setup
func init() { // runs when imported
	var err error
	addr_sending_world_view = &net.UDPAddr{
		IP:   net.ParseIP(UDP_world_view_send_ip), // Use "127.0.0.1" instead of "localhost" for consistency
		Port: UDP_world_view_send_port,
	}

	fmt.Println("DiualUDP")
	conn_sending_world_view, err = net.DialUDP("udp", nil, addr_sending_world_view)
	if err != nil {
		log.Fatalf("Failed to initialize world view send UDP connection: %v", err)
	}

	// World view receiving UDP connection setup
	addr_receiving_world_view = &net.UDPAddr{
		IP:   net.ParseIP(UDP_world_view_receive_ip),
		Port: UDP_world_view_receive_port,
	}

	fmt.Println("ListenUDP")
	conn_receiving_world_view, err = net.ListenUDP("udp", addr_receiving_world_view)
	if err != nil {
		log.Fatalf("Failed to initialize world view receive UDP connection: %v", err)
	}

}

//

func Send_elevator_world_view() {
	message := "PC1: World view message"
	_, err := conn_sending_world_view.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Sent world view from PC1")

}

func Receive_elevator_world_view(world_view_resever_chan chan string) {
	buffer := make([]byte, 1024)

	for {
		conn_receiving_world_view.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT) * time.Second)) // Setter timeout for motta adressen
		n, _, err := conn_receiving_world_view.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Failed to read from udp", err)
			time.Sleep(500 * time.Millisecond)
		} else {
			message_str := string(buffer[:n])

			world_view_resever_chan <- message_str
		}

	}
}
