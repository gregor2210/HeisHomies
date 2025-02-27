package connectivity

import (
	"Driver-go/fsm"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

const (
	//Device ID, to make it use right ports and ip. For easyer development
	//starts at 0
	ID = 0
	// Timeout for receiving UDP messages
	TIMEOUT = 3

	// World view receiving UDP connection setup

)

var (
	// What ID will given id listen to and dial to
	TCP_listen_IDs []int
	TCP_dial_IDs   []int

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
	TCP_world_view_send_ips = []string{"localhost:8080", "localhost:8070", "localhost:8060"}
	TCP_listen_conns        = [3]net.Conn{}

	other_elevatorID_order []int
)

// // World view sending TCP connection setup
func init() { // runs when imported
	//ALLE SKAL LISTENE PÅ SIN IP MEN DAILTE DE ANDRES!!!

}

func SetupTCPListen() {

}

// Serialize the struct
func SerializeElevator(wv Worldview_package) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wv)
	return buf.Bytes(), err
}

func DeserializeElevator(data []byte) (Worldview_package, error) {
	var wv Worldview_package
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wv)
	return wv, err
}

// Sender world view, i form av sin elevator struct, i form av bytes
func Send_elevator_world_view() {
	elv_struct := fsm.GetElevatorStruct()

	world_view_package_struct := New_Worldview_package(ID, elv_struct)

	elv_data, ser_err := SerializeElevator(world_view_package_struct)
	if ser_err != nil {
		fmt.Println("Serializing failed", ser_err)
	}

	_, err := conn_sending_world_view.Write(elv_data)
	if err != nil {
		fmt.Println("Failed to write", err)
	} else {
		fmt.Printf("Sent world view from PC%d\n", ID)
	}
}

func Receive_elevator_world_view_distributor(world_view_resever_chan chan Worldview_package) {
	for i := 0; i < len(conn_receiving_world_view); i++ {
		go Receive_elevator_world_view(world_view_resever_chan, conn_receiving_world_view[i], other_elevatorID_order[i])
	}
}

// Mottar verdensbilde fra andre heiser, i form av elevator structen deres
func Receive_elevator_world_view(world_view_resever_chan chan Worldview_package, conn_receiving_world_view *net.UDPConn, elevatorID_receving_from int) {
	buffer := make([]byte, 1024)

	for {
		conn_receiving_world_view.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT) * time.Second)) // Setter timeout for motta adressen
		n, _, err := conn_receiving_world_view.ReadFromUDP(buffer)
		if err != nil {
			//fmt.Println("Failed to read from udp error:", err)

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Timeout occurred while reading from UDP: %v\n", netErr)
				SetElevatorOffline(elevatorID_receving_from)

			} else {
				fmt.Printf("Error occurred while reading from UDP: %v\n", err)
				// Usikker på om denne skjer dersom ting blir corrupt på vein
			}

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

		SetElevatorOnline(elevatorID_receving_from)
		world_view_resever_chan <- elv_struct

	}
}
