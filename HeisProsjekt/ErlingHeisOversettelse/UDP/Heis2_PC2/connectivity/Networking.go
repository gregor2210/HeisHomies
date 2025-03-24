package connectivity

import (
	"Driver-go/fsm"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const (
	//Device ID, to make it use right ports and ip. For easyer development
	//starts at 0
	ID = 1
	// TimeOut for receiving UDP messages
	TimeOut = 3

	// World view receiving UDP connection setup

)

var (
	// World view sending UDP connection setup
	UDP_worldView_send_port int
	UDP_worldView_send_ip   string
	// World view sending UDP connection setup
	addr_sending_worldView     *net.UDPAddr
	conn_sending_worldView     *net.UDPConn
	UDP_worldView_receive_port []int
	UDP_worldView_receive_ip   []string
	conn_receiving_worldView   []*net.UDPConn

	other_elevatorID_order []int
)

func init() { // runs when imported. To setup the right ports and ip for the device
	if ID == 0 {
		// World view sending UDP connection setup
		UDP_worldView_send_port = 8080
		UDP_worldView_send_ip = "224.0.0.1"
		// World view receiving UDP connection setup. Multiple ports and IPs can be added
		other_elevatorID_order = []int{1, 2}
		UDP_worldView_receive_port = []int{8070, 8060}
		UDP_worldView_receive_ip = []string{"224.0.0.1", "224.0.0.1"}
		//addr_receiving_worldView *net.UDPAddr
	} else if ID == 1 {
		// World view sending UDP connection setup
		UDP_worldView_send_port = 8070
		UDP_worldView_send_ip = "224.0.0.1"
		// World view receiving UDP connection setup. Multiple ports and IPs can be added
		other_elevatorID_order = []int{0, 2}
		UDP_worldView_receive_port = []int{8080, 8060}
		UDP_worldView_receive_ip = []string{"224.0.0.1", "224.0.0.1"}
	} else if ID == 2 {
		// World view sending UDP connection setup
		UDP_worldView_send_port = 8060
		UDP_worldView_send_ip = "224.0.0.1"
		// World view receiving UDP connection setup. Multiple ports and IPs can be added
		other_elevatorID_order = []int{0, 1}
		UDP_worldView_receive_port = []int{8080, 8070}
		UDP_worldView_receive_ip = []string{"224.0.0.1", "224.0.0.1"}
	} else {
		log.Fatalf("Invalid ID: %v", ID)
	}
}

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

		// Vi måtte åpne resuadder os sol socket før listen, Use net.ListenConfig{} to set socket options before the socket is bound.
		lc := net.ListenConfig{
			Control: func(network, address string, c syscall.RawConn) error {
				return c.Control(func(fd uintptr) {
					err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
					if err != nil {
						log.Fatalf("Failed to initialize receve world view: SO_REUSEADDR failed: %v", err)
					}

					err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
					if err != nil {
						log.Fatalf("Failed to initialize receve world view: SO_REUSEPORT failed: %v", err)
					}
				})
			},
		}

		conn, err := lc.ListenPacket(context.Background(), "udp", addr.String())
		if err != nil {
			log.Fatalf("Failed to initialize world view receive UDP connection (Listen): %v", err)
		}

		conn_receiving_worldView = append(conn_receiving_worldView, conn.(*net.UDPConn))
	}

}

// Serialize the struct
func SerializeElevator(wv WorldviewPackage) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wv)
	return buf.Bytes(), err
}

func DeserializeElevator(data []byte) (WorldviewPackage, error) {
	var wv WorldviewPackage
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wv)
	return wv, err
}

// Sender world view, i form av sin elevator struct, i form av bytes
func Send_elevator_worldView() {
	elv_struct := fsm.GetElevatorStruct()

	worldView_package_struct := NewWorldviewPackage(ID, elv_struct)

	elv_data, ser_err := SerializeElevator(worldView_package_struct)
	if ser_err != nil {
		fmt.Println("Serializing failed", ser_err)
	}

	_, err := conn_sending_worldView.Write(elv_data)
	if err != nil {
		fmt.Println("Failed to write", err)
	} else {
		fmt.Printf("Sent world view from PC%d\n", ID)
	}
}

func Receive_elevator_worldView_distributor(worldView_resever_chan chan WorldviewPackage) {
	for i := 0; i < len(conn_receiving_worldView); i++ {
		go Receive_elevator_worldView(worldView_resever_chan, conn_receiving_worldView[i], other_elevatorID_order[i])
	}
}

// Mottar verdensbilde fra andre heiser, i form av elevator structen deres
func Receive_elevator_worldView(worldView_resever_chan chan WorldviewPackage, conn_receiving_worldView *net.UDPConn, elevatorID_receving_from int) {
	buffer := make([]byte, 1024)

	for {
		conn_receiving_worldView.SetReadDeadline(time.Now().Add(time.Duration(TimeOut) * time.Second)) // Setter TimeOut for motta adressen
		n, _, err := conn_receiving_worldView.ReadFromUDP(buffer)
		if err != nil {
			//fmt.Println("Failed to read from udp error:", err)

			if netErr, ok := err.(net.Error); ok && netErr.TimeOut() {
				fmt.Printf("TimeOut occurred while reading from UDP: %v\n", netErr)
				SetElevatorOffline(elevatorID_receving_from)

			} else {
				fmt.Printf("Error occurred while reading from UDP: %v\n", err)
				// Usikker på om denne skjer dersom ting blir corrupt på vein
			}

			//time.Sleep(500 * time.Millisecond)
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
		worldView_resever_chan <- elv_struct

	}
}
