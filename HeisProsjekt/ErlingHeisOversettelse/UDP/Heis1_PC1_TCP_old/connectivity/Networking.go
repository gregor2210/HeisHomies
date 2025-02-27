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
	ID = 0
	// Timeout for receiving UDP messages
	TIMEOUT = 3

	// World view receiving UDP connection setup

)

var (
	// What ID will given id listen to and dial to
	TCP_listen_ID                  = ID
	TCP_dial_ID                    int
	TCP_who_dails_this_elevator_ID int

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
	TCP_world_view_send_ips = []string{"localhost:8080", "localhost:8070", "localhost:8060"}
	TCP_listen_conns        = [3]net.Conn{}
)

// // World view sending TCP connection setup
func init() { // runs when imported
	//ALLE SKAL LISTENE PÅ SIN IP MEN DAILTE DE ANDRES!!!
	if ID == 0 {
		TCP_dial_ID = 2
		TCP_who_dails_this_elevator_ID = 1
	} else if ID == 1 {
		TCP_dial_ID = 0
		TCP_who_dails_this_elevator_ID = 2
	} else if ID == 2 {
		TCP_dial_ID = 1
		TCP_who_dails_this_elevator_ID = 0
	} else {
		fmt.Println("Invalid ID")
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

func TCP_setup(TCP_receive_channel chan Worldview_package, TCP_send_channel_listen chan Worldview_package, TCP_send_channel_dail chan Worldview_package) {
	// Start server
	go TCP_server_setup(TCP_listen_ID, TCP_receive_channel, TCP_send_channel_listen)

	// Start client
	go TCP_client_setup(TCP_dial_ID, TCP_receive_channel, TCP_send_channel_dail)
	fmt.Println("TCP setup done")
}

func TCP_server_setup(listen_ID int, TCP_receive_channel chan Worldview_package, TCP_send_channel_listen chan Worldview_package) {
	// Listen to all incoming messages
	//fmt.Println("Starting server")

	ip := TCP_world_view_send_ips[listen_ID]
	fmt.Println("Server listening on ip: ", ip)
	ln, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println("Error in TCP_server_setup")
		fmt.Println(err)
	}
	for {
		fmt.Println("Waiting for Accept")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error in TCP_server_setup")
			fmt.Println(err)
		}
		fmt.Println("Setting elevator " + fmt.Sprint(TCP_who_dails_this_elevator_ID) + " online")
		SetElevatorOnline(TCP_who_dails_this_elevator_ID)

		error_chan := make(chan struct{})

		// Launch goroutines for sending and receiving
		go handle_receive(conn, TCP_receive_channel, TCP_who_dails_this_elevator_ID, error_chan)
		go handle_send(conn, TCP_send_channel_listen, TCP_who_dails_this_elevator_ID, error_chan)

		<-error_chan
		conn.Close()
	}
}

func TCP_client_setup(dial_ID int, TCP_receive_channel chan Worldview_package, TCP_send_channel_dail chan Worldview_package) {
	ip := TCP_world_view_send_ips[dial_ID]

	for {
		fmt.Printf("Trying to dail to ip: %s\n", ip)
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			fmt.Println("Connection failed, retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("Connected to", ip)

		error_chan := make(chan struct{})

		SetElevatorOnline(dial_ID) //setting status of connected elevator to online
		// Launch goroutines for sending and receiving
		go handle_receive(conn, TCP_receive_channel, dial_ID, error_chan)
		go handle_send(conn, TCP_send_channel_dail, dial_ID, error_chan)

		<-error_chan
		conn.Close()

	}
}

func handle_receive(conn net.Conn, TCP_receive_channel chan Worldview_package, ID_of_connected_elevator int, error_chan chan struct{}) {
	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(ID_of_connected_elevator))
	for {
		select {
		case <-error_chan:
			return

		default:
			// Replace with actual receiving logic
			buffer := make([]byte, 1024)
			//conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error receiving or timedout, closing receive goroutine and conn")
				SetElevatorOffline(ID_of_connected_elevator) //setting status of connected elevator to offline
				close(error_chan)
				return
			}

			//deserialize the buffer to worldview package
			receved_world_view_package, err := DeserializeElevator(buffer)
			if err != nil {
				log.Fatal("failed to deserialize:", err)
			}

			TCP_receive_channel <- receved_world_view_package

		}

	}
}

func handle_send(conn net.Conn, TCP_send_channel chan Worldview_package, ID_of_connected_elevator int, error_chan chan struct{}) {
	fmt.Println("HANDLE SEND STARTED, ID: " + fmt.Sprint(ID_of_connected_elevator))
	for {
		select {
		case <-error_chan:
			fmt.Println("Error_chan closed")
			return // Exit goroutine if error_chan is closed

		case send_world_view_package := <-TCP_send_channel:
			fmt.Println("RECEIVED WORLD VIEW PACKAGE TO SEND")
			//serialize the world view package
			serialized_world_view_package, err := SerializeElevator(send_world_view_package)
			if err != nil {
				log.Fatal("failed to serialize:", err)
			}

			_, err = conn.Write(serialized_world_view_package)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(ID_of_connected_elevator) //setting status of connected elevator to offline
				close(error_chan)                            // Notify main loop
				return
			}
		}
	}
	fmt.Println("HANDLE SEND ENDED, ID: " + fmt.Sprint(ID_of_connected_elevator))
}

func Send_elevator_world_view(TCP_send_channel_listen chan Worldview_package, TCP_send_channel_dail chan Worldview_package) {
	fmt.Println("Start of sedning worldview")
	// Send world view to all other elevators
	elv_struct := fsm.GetElevatorStruct()
	world_view_package_struct := New_Worldview_package(ID, elv_struct)
	fmt.Println("HERRRRRRRRRRRR")

	if IsOnline(TCP_who_dails_this_elevator_ID) {
		TCP_send_channel_dail <- world_view_package_struct
		fmt.Println("Sending to elevator ", TCP_who_dails_this_elevator_ID)
	} else {
		fmt.Println("Elevator ", TCP_who_dails_this_elevator_ID, " is offline, ERROR 1")
	}
	if IsOnline(TCP_dial_ID) {
		TCP_send_channel_listen <- world_view_package_struct
		fmt.Println("Sending to elevator ", TCP_dial_ID)
	} else {
		fmt.Println("Elevator ", TCP_dial_ID, " is offline, ERROR 2")
	}
}
