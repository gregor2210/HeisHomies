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
	ID = 2
	// Timeout for receiving UDP messages
	TIMEOUT = 3

	// World view receiving UDP connection setup

)

var (
	// World view sending UDP connection setup
	UDP_world_view_send_port int
	UDP_world_view_send_ip   string
	// World view sending UDP connection setup
	addr_sending_world_view     *net.UDPAddr
	conn_sending_world_view     *net.UDPConn
	UDP_world_view_receive_port []int
	UDP_world_view_receive_ip   []string
	conn_receiving_world_view   []*net.UDPConn

	other_elevatorID_order []int
)

func init() { // runs when imported. To setup the right ports and ip for the device
	if ID == 0 {
		// World view sending UDP connection setup
		UDP_world_view_send_port = 8080
		UDP_world_view_send_ip = "224.0.0.1"
		// World view receiving UDP connection setup. Multiple ports and IPs can be added
		other_elevatorID_order = []int{1, 2}
		UDP_world_view_receive_port = []int{8070, 8060}
		UDP_world_view_receive_ip = []string{"224.0.0.1", "224.0.0.1"}
		//addr_receiving_world_view *net.UDPAddr
	} else if ID == 1 {
		// World view sending UDP connection setup
		UDP_world_view_send_port = 8070
		UDP_world_view_send_ip = "224.0.0.1"
		// World view receiving UDP connection setup. Multiple ports and IPs can be added
		other_elevatorID_order = []int{0, 2}
		UDP_world_view_receive_port = []int{8080, 8060}
		UDP_world_view_receive_ip = []string{"224.0.0.1", "224.0.0.1"}
	} else if ID == 2 {
		// World view sending UDP connection setup
		UDP_world_view_send_port = 8060
		UDP_world_view_send_ip = "224.0.0.1"
		// World view receiving UDP connection setup. Multiple ports and IPs can be added
		other_elevatorID_order = []int{0, 1}
		UDP_world_view_receive_port = []int{8080, 8070}
		UDP_world_view_receive_ip = []string{"224.0.0.1", "224.0.0.1"}
	} else {
		log.Fatalf("Invalid ID: %v", ID)
	}
}

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

		conn_receiving_world_view = append(conn_receiving_world_view, conn.(*net.UDPConn))
	}

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
		world_view_resever_chan <- elv_struct

	}
}
