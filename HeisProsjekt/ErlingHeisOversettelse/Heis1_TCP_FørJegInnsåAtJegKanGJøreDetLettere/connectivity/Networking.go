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

const (
//Device ID, to make it use right ports and ip. For easyer development
//starts at 0

)

var (

	// What ID will given id listen to and dial to
	//row listens to column. column dailer til row

	tcpWorldViewSendIpsMatrix [NumElevators - 1][NumElevators]string
	listenDialConnMatrix      [NumElevators - 1][NumElevators]net.Conn
	tryingToSetupMatrix       = [NumElevators - 1][NumElevators]bool{}
	rescever_running_matrix   = [NumElevators - 1][NumElevators]bool{}

	//these are gonna look like the matrixes above

	// Mutex for the matrixes
	muWorldViewSendIPMatrix    sync.Mutex
	muListenDialConnMatrix     sync.Mutex
	muTryingToSetupMatrix      sync.Mutex
	mu_rescever_running_matrix sync.Mutex

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
)

// // World view sending TCP connection setup
func init() { // runs when imported

	//flag.IntVar(&ID, "id", ID, "Spesefy the id with -id")
	//flag.Parse()

	for i := 0; i < NumElevators-1; i++ {
		for j := i + 1; j < NumElevators; j++ {
			ip := "localhost:80" + fmt.Sprint(i) + fmt.Sprint(j)
			set_tcpWorldViewSendIpsMatrix(i, j, ip)
		}
	}
}

func get_tcpWorldViewSendIpsMatrix(i int, j int) string {
	muWorldViewSendIPMatrix.Lock()
	defer muWorldViewSendIPMatrix.Unlock()
	return tcpWorldViewSendIpsMatrix[i][j]
}

func set_tcpWorldViewSendIpsMatrix(i int, j int, ip string) {
	muWorldViewSendIPMatrix.Lock()
	defer muWorldViewSendIPMatrix.Unlock()
	tcpWorldViewSendIpsMatrix[i][j] = ip
}

func getListenDialConnMatrix(i int, j int) net.Conn {
	muListenDialConnMatrix.Lock()
	defer muListenDialConnMatrix.Unlock()
	return listenDialConnMatrix[i][j]
}

func setListenDialConnMatrix(i int, j int, conn net.Conn) {
	muListenDialConnMatrix.Lock()
	defer muListenDialConnMatrix.Unlock()
	listenDialConnMatrix[i][j] = conn
}

func getTryingToSetupMatrix(i int, j int) bool {
	muTryingToSetupMatrix.Lock()
	defer muTryingToSetupMatrix.Unlock()
	return tryingToSetupMatrix[i][j]
}

func setTryingToSetupMatrix(i int, j int, b bool) {
	muTryingToSetupMatrix.Lock()
	defer muTryingToSetupMatrix.Unlock()
	tryingToSetupMatrix[i][j] = b
}

func get_rescever_running_matrix(i int, j int) bool {
	mu_rescever_running_matrix.Lock()
	defer mu_rescever_running_matrix.Unlock()
	return rescever_running_matrix[i][j]
}

func set_rescever_running_matrix(i int, j int, b bool) {
	mu_rescever_running_matrix.Lock()
	defer mu_rescever_running_matrix.Unlock()
	rescever_running_matrix[i][j] = b
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

func TcpReceivingSetup(tcpReceiveChannel chan WorldviewPackage) {
	fmt.Println("Starting TCP receving setup")
	for { //for loop to keep the function running
		//Server setup
		for j := ID + 1; j < NumElevators; j++ {
			//serverConn_setup
			if !getTryingToSetupMatrix(ID, j) && !IsOnline(j) {
				fmt.Println("Starting up tcpServerSetup. ", ID, "listening for: ", j)
				go tcpServerSetup(j)
			}

			//Server rescever setup
			if IsOnline(j) && !get_rescever_running_matrix(ID, j) {
				fmt.Println("Starting handleReceive for connected elevator")
				go handleReceive(getListenDialConnMatrix(ID, j), tcpReceiveChannel, j, ID, j)
			}
		}

		//Client setup
		for i := 0; i < ID; i++ {
			//clientConn_setup
			if !getTryingToSetupMatrix(i, ID) && !IsOnline(i) {
				fmt.Println("Starting up tcpClientSetup. ", ID, "dialing to: ", i)
				go tcpClientSetup(i)
			}

			//Client rescever setup
			if IsOnline(i) && !get_rescever_running_matrix(i, ID) {
				fmt.Println("Starting handleReceive for elevator we are connected to")
				go handleReceive(getListenDialConnMatrix(i, ID), tcpReceiveChannel, i, i, ID)
			}

		}
		time.Sleep(2 * time.Second)

	}

}

func tcpServerSetup(incomingElevID int) {
	//Setting up serfor for ID to listen to incomingElevID
	setTryingToSetupMatrix(ID, incomingElevID, true)

	serverIP := get_tcpWorldViewSendIpsMatrix(ID, incomingElevID)

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

	//Set no delay til true

	fmt.Println("Elevator ", incomingElevID, " connected to elevator ", ID, ". Setting ", incomingElevID, " to online")

	setListenDialConnMatrix(ID, incomingElevID, conn)

	SetElevatorOnline(incomingElevID)

	setTryingToSetupMatrix(ID, incomingElevID, false)

	defer ln.Close()
}

func tcpClientSetup(elevDialingToID int) {
	setTryingToSetupMatrix(elevDialingToID, ID, true)

	clientIP := get_tcpWorldViewSendIpsMatrix(elevDialingToID, ID)

	for {
		fmt.Printf("Trying to dail to ip: %s\n", clientIP)
		conn, err := net.Dial("tcp", clientIP)
		if err != nil {
			fmt.Println("Dailing id ", elevDialingToID, "failed, retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("Connected to", clientIP)

		setListenDialConnMatrix(elevDialingToID, ID, conn)

		SetElevatorOnline(elevDialingToID) //setting status of connected elevator to online

		break
	}
	setTryingToSetupMatrix(elevDialingToID, ID, false)
}

func handleReceive(conn net.Conn, tcpReceiveChannel chan WorldviewPackage, connectedElevatorID int, i int, j int) {
	defer conn.Close()
	set_rescever_running_matrix(i, j, true)

	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(connectedElevatorID))
	for {
		// Replace with actual receiving logic

		err := conn.SetReadDeadline(time.Now().Add(TimeOut * time.Second))
		if err != nil {
			fmt.Println("Conn not open")
			set_rescever_running_matrix(i, j, false)
			return
		}
		var packetLength uint32
		err = binary.Read(conn, binary.BigEndian, &packetLength)
		if err != nil {
			fmt.Println("failed to read packetLength:", err)
			set_rescever_running_matrix(i, j, false)
			return
		}
		buffer := make([]byte, packetLength)
		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving or timedout, closing receive goroutine and conn")
			SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
			set_rescever_running_matrix(i, j, false)
			return
		}
		//fmt.Println("DATA MOTATT! ")

		// Remove padding before deserializing
		//trimmedData := bytes.TrimRight(buffer, "\x00")

		//fmt.Printf("trimmedData: %x\n", trimmedData)

		//deserialize the buffer to worldview package
		receivedWorldViewPackage, err := DeserializeElevator(buffer)
		if err != nil {
			log.Fatal("failed to deserialize:", err)
		}

		//Store backup worldview from incomming elevator
		StoreWorldview(receivedWorldViewPackage.ElevatorID, receivedWorldViewPackage)

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

	// Sending to all online that are connected to us
	for connectedElevatorID := ID + 1; connectedElevatorID < NumElevators; connectedElevatorID++ {
		if IsOnline(connectedElevatorID) {

			//sending first packetLength, before actual packet. Preventing packet stacking
			serverConn := getListenDialConnMatrix(ID, connectedElevatorID)
			err = binary.Write(serverConn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost.")
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
			}

			//writing acctual package
			_, err = serverConn.Write(serializedWorldViewPackage)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
			} else {
				//fmt.Println("Succes sending worldview")
			}
		}
	}
	for connectedElevatorID := 0; connectedElevatorID < ID; connectedElevatorID++ {
		if IsOnline(connectedElevatorID) {
			//sending first packetLength, before actual packet. Preventing packet stacking
			clientConn := getListenDialConnMatrix(connectedElevatorID, ID)
			err = binary.Write(clientConn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost.")
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
			}

			//writing acctual package
			_, err = clientConn.Write(serializedWorldViewPackage)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connectedElevatorID) //setting status of connected elevator to offline
			}
		}
	}
}

func SendOrderToSpecificElevator(receiverElev int, order elevio.ButtonEvent) bool {
	fmt.Println("Inside send order to spesific elevator!")
	//find correct conn
	var conn net.Conn
	if IsOnline(receiverElev) {

		if getListenDialConnMatrix(ID, receiverElev) != nil {
			conn = getListenDialConnMatrix(ID, receiverElev)
			fmt.Println("Here1")

		} else if getListenDialConnMatrix(receiverElev, ID) != nil {
			conn = getListenDialConnMatrix(receiverElev, ID)
			fmt.Println("Here2")
		} else {
			fmt.Println("No valid conn to send ORDER")
			fmt.Println("Here3")
			return false
		}
	} else {
		fmt.Println("Here4")
		return false
	}
	fmt.Println("Here5")
	//Try to send order to that conn
	a := fsm.GetElevatorStruct()
	fmt.Println("After get elv struct")

	_ = NewWorldviewPackage(ID, a)
	fmt.Println("After get WV")

	SendWorldviewPackage := NewWorldviewPackage(ID, fsm.GetElevatorStruct())
	SendWorldviewPackage.OrderBool = true
	SendWorldviewPackage.Order = order

	fmt.Println("Here6")
	serializedWorldViewPackage, err := SerializeElevator(SendWorldviewPackage)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}
	fmt.Println("Here7")

	// Set the write deadline for both write operations (2 seconds)
	err = conn.SetWriteDeadline(time.Now().Add(TimeOut * time.Second))
	if err != nil {
		fmt.Println("Failed to set write deadline:", err)
		return false
	}
	fmt.Println("Here8")

	//Finding package length
	packetLength := uint32(len(serializedWorldViewPackage)) //uint32 is 4 bytes

	err = binary.Write(conn, binary.BigEndian, packetLength)
	if err != nil {
		fmt.Println("Error sending packetlength for ORDER to connected elevator, connection lost or timedout.")
		return false
	}
	fmt.Println("Success sending ORDER packetlength")

	//writing acctual package
	_, err = conn.Write(serializedWorldViewPackage)
	if err != nil {
		fmt.Println("Error sending ORDER, connection lost.  or timedout")
		return false
	}
	//disable SetwriteDeadline
	conn.SetWriteDeadline(time.Time{})
	fmt.Println("Success sending ORDER")

	//everyting worked!
	return true
}
