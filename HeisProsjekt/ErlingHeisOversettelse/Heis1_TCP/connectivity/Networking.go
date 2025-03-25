package connectivity

//https://pkg.go.dev/net#KeepAliveConfig
import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	// Matrix of TCP IPs where rows are listeners and columns are dialers
	tcpWorldViewSendIPMatrix [NumElevators - 1][NumElevators]string

	// Matrix of active TCP connections between elevators
	listenDialConnMatrix [NumElevators - 1][NumElevators]net.Conn

	// Tracks active connection attempts between elevators
	tryingToSetupMatrix = [NumElevators - 1][NumElevators]bool{}

	// Matrix indicating active receiver goroutines
	receiverRunningMatrix = [NumElevators - 1][NumElevators]bool{}

	// Mutex for the matrixes
	muWorldViewSendIPMatrix sync.Mutex
	muListenDialConnMatrix  sync.Mutex
	muTryingToSetupMatrix   sync.Mutex
	muReceiverRunningMatrix sync.Mutex
)

func init() { // runs when imported

	// This function setup the tcpWorldViewSendIPMatrix
	// If UseIPs is true, use network IPs

	if UseIPs {
		if NumElevators > len(IPs) {
			log.Fatal("NumElevators larger then amount of IPs")
		}

		for server := 0; server < NumElevators-1; server++ {
			for client := server + 1; client < NumElevators; client++ {
				ip := IPs[server] + ":80" + fmt.Sprint(server) + fmt.Sprint(client)
				setTcpWorldViewSendIPMatrix(server, client, ip)
			}
		}

	} else {

		// If UseIPs is false, set IPs to localhost for local testing
		for server := 0; server < NumElevators-1; server++ {
			for client := server + 1; client < NumElevators; client++ {
				ip := "localhost:80" + fmt.Sprint(server) + fmt.Sprint(client)
				setTcpWorldViewSendIPMatrix(server, client, ip)
			}
		}
	}
}

// Thread-safe access functions for the connection matrices
func getTcpWorldViewSendIPMatrix(server int, client int) string {
	muWorldViewSendIPMatrix.Lock()
	defer muWorldViewSendIPMatrix.Unlock()
	return tcpWorldViewSendIPMatrix[server][client]
}

func setTcpWorldViewSendIPMatrix(server int, client int, ip string) {
	muWorldViewSendIPMatrix.Lock()
	defer muWorldViewSendIPMatrix.Unlock()
	tcpWorldViewSendIPMatrix[server][client] = ip
}

func getListenDialConnMatrix(server int, client int) net.Conn {
	muListenDialConnMatrix.Lock()
	defer muListenDialConnMatrix.Unlock()
	return listenDialConnMatrix[server][client]
}

func setListenDialConnMatrix(server int, client int, conn net.Conn) {
	muListenDialConnMatrix.Lock()
	defer muListenDialConnMatrix.Unlock()
	listenDialConnMatrix[server][client] = conn
}

func getTryingToSetupMatrix(server int, client int) bool {
	muTryingToSetupMatrix.Lock()
	defer muTryingToSetupMatrix.Unlock()
	return tryingToSetupMatrix[server][client]
}

func setTryingToSetupMatrix(server int, client int, b bool) {
	muTryingToSetupMatrix.Lock()
	defer muTryingToSetupMatrix.Unlock()
	tryingToSetupMatrix[server][client] = b
}

func getReceiverRunningMatrix(server int, client int) bool {
	muReceiverRunningMatrix.Lock()
	defer muReceiverRunningMatrix.Unlock()
	return receiverRunningMatrix[server][client]
}

func setReceiverRunningMatrix(server int, client int, b bool) {
	muReceiverRunningMatrix.Lock()
	defer muReceiverRunningMatrix.Unlock()
	receiverRunningMatrix[server][client] = b
}

// Serialize the struct
func serializeElevator(wv WorldviewPackage) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wv)
	return buf.Bytes(), err
}

// Deserialize the incomming bytes
func deserializeElevator(data []byte) (WorldviewPackage, error) {
	var wv WorldviewPackage
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wv)
	return wv, err
}

// Function that closes all connections and sets all elevators offline
func CloseAllConnections() {
	for server := 0; server < NumElevators-1; server++ {
		for client := server + 1; client < NumElevators; client++ {
			if getListenDialConnMatrix(server, client) != nil {
				getListenDialConnMatrix(server, client).Close()
			}
			// Everything will close when the read/write fails
		}
	}
}

func TcpReceivingSetup(tcpReceiveChannel chan WorldviewPackage) {
	// Every 2 seconds, attempts to:
	// 1. Set up TCP servers and clients
	// 2. Establish missing connections
	// 3. Start receiver goroutines for active connections

	loopTimer := 2

	for {

		if IsSelfOnline() {
			// Server setup
			for client := ID + 1; client < NumElevators; client++ {
				if !getTryingToSetupMatrix(ID, client) && !IsOnline(client) && !getReceiverRunningMatrix(ID, client) {
					fmt.Println("Starting up tcpServerSetup. ", ID, "listening for: ", client)
					go tcpServerSetup(client)
				}

				// Server receiver setup
				if IsOnline(client) && !getReceiverRunningMatrix(ID, client) {
					fmt.Println("Starting handleReceive for connected elevator")
					go handleReceive(getListenDialConnMatrix(ID, client), tcpReceiveChannel, client, ID, client)
				}
			}

			//Client setup
			for server := 0; server < ID; server++ {
				//Client connection setup
				if !getTryingToSetupMatrix(server, ID) && !IsOnline(server) && !getReceiverRunningMatrix(server, ID) {
					fmt.Println("Starting up tcpClientSetup. ", ID, "dialing to: ", server)
					go tcpClientSetup(server)
				}

				//Client rescever setup
				if IsOnline(server) && !getReceiverRunningMatrix(server, ID) {
					fmt.Println("Starting handleReceive for elevator we are connected to")
					go handleReceive(getListenDialConnMatrix(server, ID), tcpReceiveChannel, server, server, ID)
				}

			}
		}
		time.Sleep(time.Duration(loopTimer) * time.Second)

	}

}

// Setting up server for self (ID) to listen to elevator (incomingElevID)
func tcpServerSetup(incomingElevID int) {

	setTryingToSetupMatrix(ID, incomingElevID, true)

	serverIP := getTcpWorldViewSendIPMatrix(ID, incomingElevID)

	fmt.Println("Server listening on ip: ", serverIP)
	ln, err := net.Listen("tcp", serverIP)
	if err != nil {
		fmt.Println("Error in tcpServerSetup", serverIP, err)
	}

	fmt.Println("Waiting for Accept:", serverIP)
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error in tcpServerSetup:", serverIP, err)
	}

	fmt.Println("Elevator ", incomingElevID, " connected to elevator ", ID, ". Setting ", incomingElevID, " to online")

	setListenDialConnMatrix(ID, incomingElevID, conn)

	SetElevatorOnline(incomingElevID)

	setTryingToSetupMatrix(ID, incomingElevID, false)

	defer ln.Close()
}

// Setting up client for self (ID) to dial a server (elevDialingToID)
func tcpClientSetup(elevDialingToID int) {

	setTryingToSetupMatrix(elevDialingToID, ID, true)

	clientIP := getTcpWorldViewSendIPMatrix(elevDialingToID, ID)

	for {

		// Will try to dial every second until it connects
		fmt.Printf("Trying to dail to ip: %s\n", clientIP)
		conn, err := net.DialTimeout("tcp", clientIP, TimeOut*time.Second) // Timeout after 2 TimeOut (config) seconds
		if err != nil {
			fmt.Println("Dailing id ", elevDialingToID, "failed, retrying in 1/2 seconds...")
			time.Sleep(500 * time.Millisecond)
			continue
		}

		fmt.Println("Connected to", clientIP)

		setListenDialConnMatrix(elevDialingToID, ID, conn)

		SetElevatorOnline(elevDialingToID)

		break
	}
	setTryingToSetupMatrix(elevDialingToID, ID, false)
}

func handleReceive(conn net.Conn, tcpReceiveChannel chan WorldviewPackage, connectedElevatorID int, server int, client int) {
	// Runs as a goroutine to handle incoming messages:
	// 1. Sets a read deadline. If no data is received within TimeOut, the connection is closed.
	// 2. Reads the length of the incoming packet.
	// 3. Reads the packet of the specified length.
	// 4. Deserializes and sends the worldview package to tcpReceiveChannel.

	defer conn.Close()
	setReceiverRunningMatrix(server, client, true)

	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(connectedElevatorID))
	for {
		// Replace with actual receiving logic

		// Setting read deadline
		err := conn.SetReadDeadline(time.Now().Add(TimeOut * time.Second))
		if err != nil {
			fmt.Println("Conn not open:", err)
			SetElevatorOffline(connectedElevatorID)
			setReceiverRunningMatrix(server, client, false)
			return
		}

		// Read packet length
		var packetLength uint32
		err = binary.Read(conn, binary.BigEndian, &packetLength)
		if err != nil {
			SetElevatorOffline(connectedElevatorID)
			fmt.Println("Failed to read packetLength:", err)
			setReceiverRunningMatrix(server, client, false)
			return
		}

		// Read incomming worldview packet (bytes)
		buffer := make([]byte, packetLength)
		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving or timedout, closing receive goroutine and conn: ", err)
			SetElevatorOffline(connectedElevatorID)
			setReceiverRunningMatrix(server, client, false)
			return
		}

		// Deserialize the buffer to worldview package
		receivedWorldViewPackage, err := deserializeElevator(buffer)
		if err != nil {
			fmt.Println("failed to deserialize:", err)
			fmt.Println("Seting elevator offline")
			SetElevatorOffline(connectedElevatorID)
			setReceiverRunningMatrix(server, client, false)
			return
		}

		// Store backup worldview from incomming elevator
		StoreWorldview(receivedWorldViewPackage.ElevatorID, receivedWorldViewPackage)

		tcpReceiveChannel <- receivedWorldViewPackage
	}

}

func SendWorldView() {
	// Tries to send a worldview package to all online elevators
	// First to those connected to us (servers), then to those we are connected to (clients)
	SendWorldviewPackage := NewWorldviewPackage(ID, fsm.GetElevatorStruct())
	serializedWorldViewPackage, err := serializeElevator(SendWorldviewPackage)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}

	if len(serializedWorldViewPackage) > MaxPacketSize {
		log.Fatal("error: serialized data too large")
	}

	//Finding package length
	packetLength := uint32(len(serializedWorldViewPackage)) //uint32 is 4 bytes

	// Send to elevators where we are server
	for connectedElevatorID := ID + 1; connectedElevatorID < NumElevators; connectedElevatorID++ {
		if IsOnline(connectedElevatorID) {

			// Send packet length first to avoid message stacking
			serverConn := getListenDialConnMatrix(ID, connectedElevatorID)

			// Start write session with TimeOut
			err = serverConn.SetWriteDeadline(time.Now().Add(TimeOut * time.Second))
			if err != nil {
				fmt.Println("Failed to set write deadline for server write:", err, connectedElevatorID)
				SetElevatorOffline(connectedElevatorID)
				serverConn.Close()
				continue

			}
			// Write packet length
			err = binary.Write(serverConn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost: ", err)
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
				serverConn.Close()
				continue
			}

			// Write actual package
			_, err = serverConn.Write(serializedWorldViewPackage)
			if err != nil {
				fmt.Println("Error sending, connection lost: ", err)
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
				fmt.Println("Elevator was set to ofline")
				serverConn.Close()
				continue
			}

			// Disable write deadline after transmission
			serverConn.SetWriteDeadline(time.Time{})
		}
	}

	// Send to elevators where we are client
	for connectedElevatorID := 0; connectedElevatorID < ID; connectedElevatorID++ {
		if IsOnline(connectedElevatorID) {

			clientConn := getListenDialConnMatrix(connectedElevatorID, ID)

			err = clientConn.SetWriteDeadline(time.Now().Add(TimeOut * time.Second))
			if err != nil {
				fmt.Println("Failed to set write deadline for client write:", err, connectedElevatorID)
				SetElevatorOffline(connectedElevatorID)
				clientConn.Close()
				continue
			}

			err = binary.Write(clientConn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost: ", err)
				SetElevatorOffline(connectedElevatorID)
				clientConn.Close()
				continue
			}

			_, err = clientConn.Write(serializedWorldViewPackage)
			if err != nil {
				fmt.Println("Error sending, connection lost: ", err)
				SetElevatorOffline(connectedElevatorID)
				clientConn.Close()
				continue
			}

			clientConn.SetWriteDeadline(time.Time{})
		}
	}
}

func SendOrderToSpecificElevator(receiverElev int, order elevio.ButtonEvent) bool {
	// Find the correct connection, regardless of whether this elevator is dialing or listening
	var conn net.Conn
	if IsOnline(receiverElev) {
		if ID < (NumElevators-1) && getListenDialConnMatrix(ID, receiverElev) != nil {
			conn = getListenDialConnMatrix(ID, receiverElev)

		} else if receiverElev < (NumElevators-1) && getListenDialConnMatrix(receiverElev, ID) != nil {
			conn = getListenDialConnMatrix(receiverElev, ID)
		} else {
			fmt.Println("No valid conn to send ORDER")
			return false
		}
	} else {
		return false
	}

	SendWorldviewPackage := NewWorldviewPackage(ID, fsm.GetElevatorStruct())
	SendWorldviewPackage.OrderBool = true
	SendWorldviewPackage.Order = order

	serializedWorldViewPackage, err := serializeElevator(SendWorldviewPackage)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}

	// Set the write deadline for both write operations (2 seconds)
	err = conn.SetWriteDeadline(time.Now().Add(TimeOut * time.Second))
	if err != nil {
		fmt.Println("Failed to set write deadline:", err)
		return false
	}

	// Finding package length
	packetLength := uint32(len(serializedWorldViewPackage)) // uint32 is 4 bytes

	err = binary.Write(conn, binary.BigEndian, packetLength)
	if err != nil {
		fmt.Println("Error sending packetlength for ORDER to connected elevator, connection lost or timedout: ", err)
		SetElevatorOffline(receiverElev)
		conn.Close()
		return false
	}
	fmt.Println("Success sending ORDER packetlength")

	// Writing actual package
	_, err = conn.Write(serializedWorldViewPackage)
	if err != nil {
		fmt.Println("Error sending ORDER, connection lost.  or timedout: ", err)
		SetElevatorOffline(receiverElev)
		conn.Close()
		return false
	}

	// Disable SetwriteDeadline
	conn.SetWriteDeadline(time.Time{})

	return true
}
