package connectivity

import (
	"Driver-go/fsm"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	// TimeOut for receiving UDP messages
	TimeOut = 3
	// World view sending UDP connection setup
	UDP_worldView_send_port = 8080
	UDP_worldView_send_ip   = "127.0.0.1"

	// World view receiving UDP connection setup

)

var (
	// World view sending UDP connection setup
	addr_sending_worldView *net.UDPAddr
	conn_sending_worldView *net.UDPConn

	// World view receiving UDP connection setup. Multiple ports and IPs can be added
	UDP_worldView_receive_port = []int{8070}
	UDP_worldView_receive_ip   = []string{"127.0.0.1"}
	//addr_receiving_worldView *net.UDPAddr
	conn_receiving_worldView []*net.UDPConn
)

// // World view sending UDP connection setup
func init() { // runs when imported
	var err error
	addr_sending_worldView = &net.UDPAddr{
		IP:   net.ParseIP(UDP_worldView_send_ip), // Use "127.0.0.1" instead of "localhost" for consistency
		Port: UDP_worldView_send_port,
	}

	fmt.Println("DiualUDP")
	conn_sending_worldView, err = net.DialUDP("udp", nil, addr_sending_worldView)
	if err != nil {
		log.Fatalf("Failed to initialize world view send UDP connection: %v", err)
	}

	// World view receiving UDP connection setup
	for i := 0; i < len(UDP_worldView_receive_port); i++ {
		addr := &net.UDPAddr{
			IP:   net.ParseIP(UDP_worldView_receive_ip[i]),
			Port: UDP_worldView_receive_port[i],
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatalf("Failed to initialize world view receive UDP connection: %v", err)
		}
		conn_receiving_worldView = append(conn_receiving_worldView, conn)
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
func Send_elevator_worldView() {
	elv_struct := fsm.GetElevatorStruct()

	elv_data, ser_err := SerializeElevator(elv_struct)
	if ser_err != nil {
		fmt.Println("Serializing failed", ser_err)
	}

	_, err := conn_sending_worldView.Write(elv_data)
	if err != nil {
		fmt.Println("Failed to write", err)
	} else {
		fmt.Println("Sent world view from PC1")
	}
}

func Receive_elevator_worldView_distributor(worldView_resever_chan chan fsm.Elevator) {
	for i := 0; i < len(conn_receiving_worldView); i++ {
		go Receive_elevator_worldView(worldView_resever_chan, conn_receiving_worldView[i])
	}
}

// Mottar verdensbilde fra andre heiser, i form av elevator structen deres
func Receive_elevator_worldView(worldView_resever_chan chan fsm.Elevator, conn_receiving_worldView *net.UDPConn) {
	buffer := make([]byte, 1024)

	for {
		conn_receiving_worldView.SetReadDeadline(time.Now().Add(time.Duration(TimeOut) * time.Second)) // Setter TimeOut for motta adressen
		n, _, err := conn_receiving_worldView.ReadFromUDP(buffer)
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
		worldView_resever_chan <- elv_struct
	}
}
