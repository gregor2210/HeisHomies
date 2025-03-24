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
	// TimeOut for receiving UDP messages
	TimeOut = 3

	// World view receiving UDP connection setup

)

var (
	// What ID will given id listen to and dial to
	TCP_listen_ID                 = ID
	TCP_dial_ID                   int
	TCP_who_dails_this_ElevatorID int

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
	TCP_worldView_send_ips = []string{"localhost:8080", "localhost:8070", "localhost:8060"}
	TCP_listen_conns       = [3]net.Conn{}
)

// // World view sending TCP connection setup
func init() { // runs when imported
	//ALLE SKAL LISTENE PÅ SIN IP MEN DAILTE DE ANDRES!!!
	if ID == 0 {
		TCP_dial_ID = 2
		TCP_who_dails_this_ElevatorID = 1
	} else if ID == 1 {
		TCP_dial_ID = 0
		TCP_who_dails_this_ElevatorID = 2
	} else if ID == 2 {
		TCP_dial_ID = 1
		TCP_who_dails_this_ElevatorID = 0
	} else {
		fmt.Println("Invalid ID")
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

func TCP_setup(tcpReceiveChannel chan WorldviewPackage, TCP_send_channel_listen chan WorldviewPackage, TCP_send_channel_dail chan WorldviewPackage) {
	// Start server
	go tcpServerSetup(TCP_listen_ID, tcpReceiveChannel, TCP_send_channel_listen)

	// Start client
	go tcpClientSetup(TCP_dial_ID, tcpReceiveChannel, TCP_send_channel_dail)
	fmt.Println("TCP setup done")
}

func tcpServerSetup(listen_ID int, tcpReceiveChannel chan WorldviewPackage, TCP_send_channel_listen chan WorldviewPackage) {
	// Listen to all incoming messages
	//fmt.Println("Starting server")

	ip := TCP_worldView_send_ips[listen_ID]
	fmt.Println("Server listening on ip: ", ip)
	ln, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println("Error in tcpServerSetup")
		fmt.Println(err)
	}
	for {
		fmt.Println("Waiting for Accept")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error in tcpServerSetup")
			fmt.Println(err)
		}
		fmt.Println("Setting elevator " + fmt.Sprint(TCP_who_dails_this_ElevatorID) + " online")
		SetElevatorOnline(TCP_who_dails_this_ElevatorID)

		error_chan := make(chan struct{})

		// Launch goroutines for sending and receiving
		go handleReceive(conn, tcpReceiveChannel, TCP_who_dails_this_ElevatorID, error_chan)
		go handle_send(conn, TCP_send_channel_listen, TCP_who_dails_this_ElevatorID, error_chan)

		<-error_chan
		conn.Close()
	}
}

func tcpClientSetup(dial_ID int, tcpReceiveChannel chan WorldviewPackage, TCP_send_channel_dail chan WorldviewPackage) {
	ip := TCP_worldView_send_ips[dial_ID]

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
		go handleReceive(conn, tcpReceiveChannel, dial_ID, error_chan)
		go handle_send(conn, TCP_send_channel_dail, dial_ID, error_chan)

		<-error_chan
		conn.Close()

	}
}

func handleReceive(conn net.Conn, tcpReceiveChannel chan WorldviewPackage, connectedElevatorID int, error_chan chan struct{}) {
	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(connectedElevatorID))
	for {
		select {
		case <-error_chan:
			return

		default:
			// Replace with actual receiving logic
			buffer := make([]byte, 1024)
			//conn.SetReadDeadline(time.Now().Add(TimeOut * time.Second))
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error receiving or timedout, closing receive goroutine and conn")
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
				close(error_chan)
				return
			}

			//deserialize the buffer to worldview package
			receivedWorldViewPackage, err := DeserializeElevator(buffer)
			if err != nil {
				log.Fatal("failed to deserialize:", err)
			}

			tcpReceiveChannel <- receivedWorldViewPackage

		}

	}
}

func handle_send(conn net.Conn, TCP_send_channel chan WorldviewPackage, connectedElevatorID int, error_chan chan struct{}) {
	fmt.Println("HANDLE SEND STARTED, ID: " + fmt.Sprint(connectedElevatorID))
	for {
		select {
		case <-error_chan:
			fmt.Println("Error_chan closed")
			return // Exit goroutine if error_chan is closed

		case SendWorldviewPackage := <-TCP_send_channel:
			fmt.Println("RECEIVED WORLD VIEW PACKAGE TO SEND")
			//serialize the world view package
			serializedWorldViewPackage, err := SerializeElevator(SendWorldviewPackage)
			if err != nil {
				log.Fatal("failed to serialize:", err)
			}

			_, err = conn.Write(serializedWorldViewPackage)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
				close(error_chan)                       // Notify main loop
				return
			}
		}
	}
	fmt.Println("HANDLE SEND ENDED, ID: " + fmt.Sprint(connectedElevatorID))
}

func Send_elevator_worldView(TCP_send_channel_listen chan WorldviewPackage, TCP_send_channel_dail chan WorldviewPackage) {
	fmt.Println("Start of sedning worldview")
	// Send world view to all other elevators
	elv_struct := fsm.GetElevatorStruct()
	worldView_package_struct := NewWorldviewPackage(ID, elv_struct)
	fmt.Println("HERRRRRRRRRRRR")

	if IsOnline(TCP_who_dails_this_ElevatorID) {
		TCP_send_channel_dail <- worldView_package_struct
		fmt.Println("Sending to elevator ", TCP_who_dails_this_ElevatorID)
	} else {
		fmt.Println("Elevator ", TCP_who_dails_this_ElevatorID, " is offline, ERROR 1")
	}
	if IsOnline(TCP_dial_ID) {
		TCP_send_channel_listen <- worldView_package_struct
		fmt.Println("Sending to elevator ", TCP_dial_ID)
	} else {
		fmt.Println("Elevator ", TCP_dial_ID, " is offline, ERROR 2")
	}
}
