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
	//Device ID, to make it use right ports and ip. For easyer development
	//starts at 0
	ID = 1
	// Timeout for receiving UDP messages
	TIMEOUT = 5

	// Worldview max package size
	PACKAGE_SIZE = 1024
)

var (
	// What ID will given id listen to and dial to
	TCP_listen_ID          = ID
	TCP_e_we_connect_to_ID int
	TCP_connected_e_ID     int

	server_ip string
	client_ip string

	server_conn net.Conn
	client_conn net.Conn

	server_trying_to_setup bool = false
	client_trying_to_setup bool = false

	server_resiever_running bool = false
	client_resiever_running bool = false

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
	TCP_world_view_send_ips = []string{"localhost:8080", "localhost:8070", "localhost:8060"}
	TCP_listen_conns        = [3]net.Conn{}
)

// // World view sending TCP connection setup
func init() { // runs when imported
	//ALLE SKAL LISTENE PÅ SIN IP MEN DAILTE DE ANDRES!!!
	if ID == 0 {
		TCP_e_we_connect_to_ID = 2
		TCP_connected_e_ID = 1
	} else if ID == 1 {
		TCP_e_we_connect_to_ID = 0
		TCP_connected_e_ID = 2
	} else if ID == 2 {
		TCP_e_we_connect_to_ID = 1
		TCP_connected_e_ID = 0
	} else {
		fmt.Println("Invalid ID")
	}
	server_ip = TCP_world_view_send_ips[TCP_listen_ID]
	client_ip = TCP_world_view_send_ips[TCP_connected_e_ID]

	// Start server, listening for incoming connections
	go TCP_server_setup()

	// Start client, dialing to other elevator
	go TCP_client_setup()
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

func TCP_receving_setup(TCP_receive_channel chan Worldview_package, TCP_send_channel_listen chan Worldview_package, TCP_send_channel_dail chan Worldview_package) {
	fmt.Println("Starting TCP receving setup")
	for {
		// Cheching if tcp is setup, if not, setup
		if !server_trying_to_setup && !IsOnline(TCP_connected_e_ID) {
			go TCP_server_setup()
		}
		if !client_trying_to_setup && !IsOnline(TCP_e_we_connect_to_ID) {
			go TCP_client_setup()
		}

		if IsOnline(TCP_connected_e_ID) && !server_resiever_running {
			fmt.Println("Starting receiving from connected elevator")
			go handle_receive(server_conn, TCP_receive_channel, TCP_connected_e_ID, "server")
		} else {
			fmt.Println("No resceving started. No elevator is connected to this elevator")
		}
		if IsOnline(TCP_e_we_connect_to_ID) && !client_resiever_running {
			fmt.Println("Starting receiving from elevator we are connected to")
			go handle_receive(client_conn, TCP_receive_channel, TCP_e_we_connect_to_ID, "client")
		} else {
			fmt.Println("No resceving started. This elevator is not connected to any other elevator")
		}

		time.Sleep(2 * time.Second)
	}

}

func TCP_server_setup() {
	// Listen to all incoming messages
	//fmt.Println("Starting server")
	server_trying_to_setup = true

	fmt.Println("Server listening on ip: ", server_ip)
	ln, err := net.Listen("tcp", server_ip)
	if err != nil {
		fmt.Println("Error in TCP_server_setup")
		fmt.Println(err)
	}

	fmt.Println("Waiting for Accept")
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error in TCP_server_setup")
		fmt.Println(err)
	}
	fmt.Println("Setting elevator " + fmt.Sprint(TCP_connected_e_ID) + " online. GOT ACCEPTED")

	server_conn = conn
	SetElevatorOnline(TCP_connected_e_ID)
	server_trying_to_setup = false
	ln.Close()
}

func TCP_client_setup() {
	client_trying_to_setup = true
	for {
		fmt.Printf("Trying to dail to ip: %s\n", client_ip)
		conn, err := net.Dial("tcp", client_ip)
		if err != nil {
			fmt.Println("Connection failed, retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("Connected to", client_ip)

		client_conn = conn
		SetElevatorOnline(TCP_e_we_connect_to_ID) //setting status of connected elevator to online

		break
	}
	client_trying_to_setup = false
}

func handle_receive(conn net.Conn, TCP_receive_channel chan Worldview_package, ID_of_connected_elevator int, conn_type string) {
	defer conn.Close()
	if conn_type == "server" {
		server_resiever_running = true
	} else {
		client_resiever_running = true
	}

	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(ID_of_connected_elevator))
	for {
		// Replace with actual receiving logic
		buffer := make([]byte, 1024)

		err := conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
		if err != nil {
			fmt.Println("Conn not open")

			if conn_type == "server" {
				server_resiever_running = false
			} else {
				client_resiever_running = false
			}
			return
		}
		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving or timedout, closing receive goroutine and conn")
			SetElevatorOffline(ID_of_connected_elevator) //setting status of connected elevator to offline
			if conn_type == "server" {
				server_resiever_running = false
			} else {
				client_resiever_running = false
				//conn.Close()
			}
			return
		}
		fmt.Printf("DATA MOTATT!")

		// Remove padding before deserializing
		//trimmedData := bytes.TrimRight(buffer, "\x00")

		//fmt.Printf("trimmedData: %x\n", trimmedData)

		//deserialize the buffer to worldview package
		receved_world_view_package, err := DeserializeElevator(buffer)
		if err != nil {
			log.Fatal("failed to deserialize:", err)
		}

		TCP_receive_channel <- receved_world_view_package
	}

}

func Send_world_view() {
	send_world_view_package := New_Worldview_package(ID, fsm.GetElevatorStruct())
	serialized_world_view_package, err := SerializeElevator(send_world_view_package)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}

	if len(serialized_world_view_package) > PACKAGE_SIZE {
		log.Fatal("error: serialized data too large")
	}

	// Pad data with zeros to make it exactly 1024 bytes
	//paddedData := make([]byte, PACKAGE_SIZE)
	//copy(paddedData, serialized_world_view_package)

	if IsOnline(TCP_connected_e_ID) {
		_, err = server_conn.Write(serialized_world_view_package)
		if err != nil {
			fmt.Println("Error sending, connection lost.")
			SetElevatorOffline(TCP_connected_e_ID) //setting status of connected elevator to offline
		}
	}
	if IsOnline(TCP_e_we_connect_to_ID) {
		_, err = client_conn.Write(serialized_world_view_package)
		if err != nil {
			fmt.Println("Error sending, connection lost.")
			SetElevatorOffline(TCP_e_we_connect_to_ID) //setting status of connected elevator to offline
		}
	}
}
