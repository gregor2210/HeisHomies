package connectivity

import (
	"Driver-go/fsm"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	//Device ID, to make it use right ports and ip. For easyer development
	//starts at 0
	ID = 2
	// TimeOut for receiving UDP messages
	TimeOut = 5

	// Worldview max package size
	MaxPacketSize = 1024
)

var (
	// What ID will given id listen to and dial to
	TCP_listen_ID           = ID
	TCP_e_we_connect_to_ID  int
	TCP_connectedElevatorID int

	serverIP string
	clientIP string

	serverConn net.Conn
	clientConn net.Conn

	server_trying_to_setup bool = false
	client_trying_to_setup bool = false

	server_resiever_running bool = false
	client_resiever_running bool = false

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
	TCP_worldView_send_ips = []string{"localhost:8080", "localhost:8070", "localhost:8060"}
	TCP_listen_conns       = [3]net.Conn{}
)

// // World view sending TCP connection setup
func init() { // runs when imported
	//ALLE SKAL LISTENE PÅ SIN IP MEN DAILTE DE ANDRES!!!
	if ID == 0 {
		TCP_e_we_connect_to_ID = 2
		TCP_connectedElevatorID = 1
	} else if ID == 1 {
		TCP_e_we_connect_to_ID = 0
		TCP_connectedElevatorID = 2
	} else if ID == 2 {
		TCP_e_we_connect_to_ID = 1
		TCP_connectedElevatorID = 0
	} else {
		fmt.Println("Invalid ID")
	}
	serverIP = TCP_worldView_send_ips[TCP_listen_ID]
	clientIP = TCP_worldView_send_ips[TCP_connectedElevatorID]

	// Start server, listening for incoming connections
	go tcpServerSetup()

	// Start client, dialing to other elevator
	go tcpClientSetup()
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

func TcpReceivingSetup(tcpReceiveChannel chan WorldviewPackage, TCP_send_channel_listen chan WorldviewPackage, TCP_send_channel_dail chan WorldviewPackage) {
	fmt.Println("Starting TCP receving setup")
	for {
		// Cheching if tcp is setup, if not, setup
		if !server_trying_to_setup && !IsOnline(TCP_connectedElevatorID) {
			fmt.Println("Starting up tcpServerSetup")
			go tcpServerSetup()
		}
		if !client_trying_to_setup && !IsOnline(TCP_e_we_connect_to_ID) {
			fmt.Println("Starting up tcpClientSetup")
			go tcpClientSetup()
		}

		if IsOnline(TCP_connectedElevatorID) && !server_resiever_running {
			fmt.Println("Starting handleReceive for connected elevator")
			go handleReceive(serverConn, tcpReceiveChannel, TCP_connectedElevatorID, "server")
		} else {
			//fmt.Println("No resceving started. No elevator is connected to this elevator")
		}
		if IsOnline(TCP_e_we_connect_to_ID) && !client_resiever_running {
			fmt.Println("Starting rhandleReceive for elevator we are connected to")
			go handleReceive(clientConn, tcpReceiveChannel, TCP_e_we_connect_to_ID, "client")
		} else {
			//fmt.Println("No resceving started. This elevator is not connected to any other elevator")
		}

		time.Sleep(2 * time.Second)
	}

}

func tcpServerSetup() {
	// Listen to all incoming messages
	//fmt.Println("Starting server")
	server_trying_to_setup = true

	fmt.Println("Server listening on ip: ", serverIP)
	ln, err := net.Listen("tcp", serverIP)
	if err != nil {
		fmt.Println("Error in tcpServerSetup")
		fmt.Println(err)
	}

	fmt.Println("Waiting for Accept")
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error in tcpServerSetup")
		fmt.Println(err)
	}
	fmt.Println("Setting elevator " + fmt.Sprint(TCP_connectedElevatorID) + " online. GOT ACCEPTED")

	serverConn = conn
	SetElevatorOnline(TCP_connectedElevatorID)
	server_trying_to_setup = false
	ln.Close()
}

func tcpClientSetup() {
	client_trying_to_setup = true
	for {
		fmt.Printf("Trying to dail to ip: %s\n", clientIP)
		conn, err := net.Dial("tcp", clientIP)
		if err != nil {
			fmt.Println("Connection failed, retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("Connected to", clientIP)

		clientConn = conn
		SetElevatorOnline(TCP_e_we_connect_to_ID) //setting status of connected elevator to online

		break
	}
	client_trying_to_setup = false
}

func handleReceive(conn net.Conn, tcpReceiveChannel chan WorldviewPackage, connectedElevatorID int, conn_type string) {
	defer conn.Close()
	if conn_type == "server" {
		server_resiever_running = true
	} else {
		client_resiever_running = true
	}

	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(connectedElevatorID))
	for {
		// Replace with actual receiving logic
		buffer := make([]byte, 1024)

		err := conn.SetReadDeadline(time.Now().Add(TimeOut * time.Second))
		if err != nil {
			fmt.Println("Conn not open")

			if conn_type == "server" {
				server_resiever_running = false
			} else {
				client_resiever_running = false
			}
			return
		}
		var packetLength uint32
		err = binary.Read(conn, binary.BigEndian, &packetLength)
		if err != nil {
			fmt.Println("failed to read packetLength:", err)
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
			SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
			if conn_type == "server" {
				server_resiever_running = false
			} else {
				client_resiever_running = false
			}
			return
		}
		fmt.Printf("DATA MOTATT! ")

		// Remove padding before deserializing
		//trimmedData := bytes.TrimRight(buffer, "\x00")

		//fmt.Printf("trimmedData: %x\n", trimmedData)

		//deserialize the buffer to worldview package
		receivedWorldViewPackage, err := DeserializeElevator(buffer)
		if err != nil {
			log.Fatal("failed to deserialize:", err)
		}

		tcpReceiveChannel <- receivedWorldViewPackage
	}

}

func SendWorldView() {
	SendWorldviewPackage := NewWorldviewPackage(ID, fsm.GetElevatorStruct())
	serializedWorldViewPackage, err := SerializeElevator(SendWorldviewPackage)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}

	if len(serializedWorldViewPackage) > MaxPacketSize {
		log.Fatal("error: serialized data too large")
	}

	// Pad data with zeros to make it exactly 1024 bytes
	//paddedData := make([]byte, MaxPacketSize)
	//copy(paddedData, serializedWorldViewPackage)

	//Finding package length
	packetLength := uint32(len(serializedWorldViewPackage)) //uint32 is 4 bytes

	if IsOnline(TCP_connectedElevatorID) {
		//sending first packetLength, before actual packet. Preventing packet stacking
		err = binary.Write(serverConn, binary.BigEndian, packetLength)
		if err != nil {
			fmt.Println("Error sending packetlength to connected elevator, connection lost.")
			SetElevatorOffline(TCP_connectedElevatorID) //setting status of connected elevator to offline
		}

		//writing acctual package
		_, err = serverConn.Write(serializedWorldViewPackage)
		if err != nil {
			fmt.Println("Error sending, connection lost.")
			SetElevatorOffline(TCP_connectedElevatorID) //setting status of connected elevator to offline
		}
	}
	if IsOnline(TCP_e_we_connect_to_ID) {
		//sending first packetLength, before actual packet. Preventing packet stacking
		err = binary.Write(clientConn, binary.BigEndian, packetLength)
		if err != nil {
			fmt.Println("Error sending packetlength to connected elevator, connection lost.")
			SetElevatorOffline(TCP_e_we_connect_to_ID) //setting status of connected elevator to offline
		}

		//writing acctual package
		_, err = clientConn.Write(serializedWorldViewPackage)
		if err != nil {
			fmt.Println("Error sending, connection lost.")
			SetElevatorOffline(TCP_e_we_connect_to_ID) //setting status of connected elevator to offline
		}
	}
}
