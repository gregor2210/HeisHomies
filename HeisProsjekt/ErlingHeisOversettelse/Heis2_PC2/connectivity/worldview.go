package connectivity

import (
	"Driver-go/fsm"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

const (
	// Timeout for receiving UDP messages
	TIMEOUT = 3
	// World view sending UDP connection setup
	UDP_world_view_send_port = 8070
	UDP_world_view_send_ip   = "127.0.0.1"

	// World view receiving UDP connection setup

)

var (
	// World view sending UDP connection setup
	addr_sending_world_view *net.UDPAddr
	conn_sending_world_view *net.UDPConn

	// World view receiving UDP connection setup. Multiple ports and IPs can be added
	UDP_world_view_receive_port = []int{8080}
	UDP_world_view_receive_ip   = []string{"127.0.0.1"}
	//addr_receiving_world_view *net.UDPAddr
	conn_receiving_world_view []*net.UDPConn
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
	for i := 0; i < len(UDP_world_view_receive_port); i++ {
		addr := &net.UDPAddr{
			IP:   net.ParseIP(UDP_world_view_receive_ip[i]),
			Port: UDP_world_view_receive_port[i],
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatalf("Failed to initialize world view receive UDP connection: %v", err)
		}

		file, err := conn.File()
		if err != nil {
			log.Fatalf("failed to get file descriptor: %v", err)
		}
		//defer file.Close()

		err = syscall.SetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			log.Fatalf("failed to set SO_REUSEADDR: %v", err)
		}
		conn_receiving_world_view = append(conn_receiving_world_view, conn)
	}
}

//

// Serialize the struct
func SerializeElevator(e fsm.Elevator) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e)
	return buf.Bytes(), err
}

func DeserializeElevator(data []byte) (fsm.Elevator, error) {
	var elv fsm.Elevator
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&elv)
	return elv, err
}

// Sender world view, i form av sin elevator struct, i form av bytes
func Send_elevator_world_view() {
	elv_struct := fsm.GetElevatorStruct()

	elv_data, ser_err := SerializeElevator(elv_struct)
	if ser_err != nil {
		fmt.Println("Serializing failed", ser_err)
	}

	_, err := conn_sending_world_view.Write(elv_data)
	if err != nil {
		fmt.Println("Failed to write", err)
	} else {
		fmt.Println("Sent world view from PC2")
	}
}

func Receive_elevator_world_view_distributor(world_view_resever_chan chan fsm.Elevator) {
	for i := 0; i < len(conn_receiving_world_view); i++ {
		go Receive_elevator_world_view(world_view_resever_chan, conn_receiving_world_view[i])
	}
}

// Mottar verdensbilde fra andre heiser, i form av elevator structen deres
func Receive_elevator_world_view(world_view_resever_chan chan fsm.Elevator, conn_receiving_world_view *net.UDPConn) {
	buffer := make([]byte, 1024)

	for {
		conn_receiving_world_view.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT) * time.Second)) // Setter timeout for motta adressen
		n, _, err := conn_receiving_world_view.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Failed to read from udp", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		//fmt.Println(n)
		if n == 20 {
			fmt.Println("No data received")
			continue
		}
		elv_struct, err := DeserializeElevator(buffer[:n])
		if err != nil {
			fmt.Println("failed to deseralize", err)
			continue
		}
		world_view_resever_chan <- elv_struct
	}
}
