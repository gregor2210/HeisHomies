package connectivity

import (
	"Driver-go/fsm"
	"Driver-go/masterSlave"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	//Device ID, to make it use right ports and ip. For easyer development
	//starts at 0
	Nr_of_elevators = 3
	// Timeout for receiving UDP messages
	TIMEOUT = 5

	// Worldview max package size
	PACKAGE_SIZE = 1024
)

var (
	ID = 1

	// What ID will given id listen to and dial to
	//row listens to column. column dailer til row
	TCP_world_view_send_ips_matrix [Nr_of_elevators - 1][Nr_of_elevators]string
	listen_dail_conn_matrix        [Nr_of_elevators - 1][Nr_of_elevators]net.Conn
	trying_to_setup_matrix         = [Nr_of_elevators - 1][Nr_of_elevators]bool{}
	rescever_running_matrix        = [Nr_of_elevators - 1][Nr_of_elevators]bool{}

	// Mutex for the matrixes
	mu_world_view_send_ips_matrix sync.Mutex
	mu_listen_dail_conn_matrix    sync.Mutex
	mu_trying_to_setup_matrix     sync.Mutex
	mu_rescever_running_matrix    sync.Mutex

	// World view sending UDP connection setup
	// Elevator (0-1, 1-2, 2-1), first is dialing, second is listening
	TCP_world_view_send_ips = []string{"localhost:8080", "localhost:8070", "localhost:8060"}
	TCP_listen_conns        = [3]net.Conn{}
)

// // World view sending TCP connection setup
func init() { // runs when imported

	flag.IntVar(&ID, "id", ID, "Spesefy the id with -id")
	flag.Parse()

	for i := 0; i < Nr_of_elevators-1; i++ {
		for j := i + 1; j < Nr_of_elevators; j++ {
			ip := "localhost:80" + fmt.Sprint(i) + fmt.Sprint(j)
			set_TCP_world_view_send_ips_matrix(i, j, ip)
		}
	}
}

func get_TCP_world_view_send_ips_matrix(i int, j int) string {
	mu_world_view_send_ips_matrix.Lock()
	defer mu_world_view_send_ips_matrix.Unlock()
	return TCP_world_view_send_ips_matrix[i][j]
}

func set_TCP_world_view_send_ips_matrix(i int, j int, ip string) {
	mu_world_view_send_ips_matrix.Lock()
	TCP_world_view_send_ips_matrix[i][j] = ip
	mu_world_view_send_ips_matrix.Unlock()
}

func get_listen_dail_conn_matrix(i int, j int) net.Conn {
	mu_listen_dail_conn_matrix.Lock()
	defer mu_listen_dail_conn_matrix.Unlock()
	return listen_dail_conn_matrix[i][j]
}

func set_listen_dail_conn_matrix(i int, j int, conn net.Conn) {
	mu_listen_dail_conn_matrix.Lock()
	listen_dail_conn_matrix[i][j] = conn
	mu_listen_dail_conn_matrix.Unlock()
}

func get_trying_to_setup_matrix(i int, j int) bool {
	mu_trying_to_setup_matrix.Lock()
	defer mu_trying_to_setup_matrix.Unlock()
	return trying_to_setup_matrix[i][j]
}

func set_trying_to_setup_matrix(i int, j int, b bool) {
	mu_trying_to_setup_matrix.Lock()
	trying_to_setup_matrix[i][j] = b
	mu_trying_to_setup_matrix.Unlock()
}

func get_rescever_running_matrix(i int, j int) bool {
	mu_rescever_running_matrix.Lock()
	defer mu_rescever_running_matrix.Unlock()
	return rescever_running_matrix[i][j]
}

func set_rescever_running_matrix(i int, j int, b bool) {
	mu_rescever_running_matrix.Lock()
	rescever_running_matrix[i][j] = b
	mu_rescever_running_matrix.Unlock()
}

// Serialize the struct
func SerializeElevator(wv masterSlave.Worldview_package) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wv)
	return buf.Bytes(), err
}

func DeserializeElevator(data []byte) (masterSlave.Worldview_package, error) {
	var wv masterSlave.Worldview_package
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wv)
	return wv, err
}

func TCP_receving_setup(TCP_receive_channel chan masterSlave.Worldview_package) {
	fmt.Println("Starting TCP receving setup")
	for { //for loop to keep the function running
		//Server setup
		for j := ID + 1; j < Nr_of_elevators; j++ {
			//Server_conn_setup
			if !get_trying_to_setup_matrix(ID, j) && !IsOnline(j) {
				fmt.Println("Starting up TCP_server_setup. ", ID, "listening for: ", j)
				go TCP_server_setup(j)
			}

			//Server rescever setup
			if IsOnline(j) && !get_rescever_running_matrix(ID, j) {
				fmt.Println("Starting handle_receive for connected elevator")
				go handle_receive(get_listen_dail_conn_matrix(ID, j), TCP_receive_channel, j, ID, j)
			}
		}

		//Client setup
		for i := 0; i < ID; i++ {
			//Client_conn_setup
			if !get_trying_to_setup_matrix(i, ID) && !IsOnline(i) {
				fmt.Println("Starting up TCP_client_setup. ", ID, "dialing to: ", i)
				go TCP_client_setup(i)
			}

			//Client rescever setup
			if IsOnline(i) && !get_rescever_running_matrix(i, ID) {
				fmt.Println("Starting handle_receive for elevator we are connected to")
				go handle_receive(get_listen_dail_conn_matrix(i, ID), TCP_receive_channel, i, i, ID)
			}

		}
		time.Sleep(2 * time.Second)

	}

}

func TCP_server_setup(incoming_e_ID int) {
	//Setting up serfor for ID to listen to incoming_e_ID
	set_trying_to_setup_matrix(ID, incoming_e_ID, true)

	server_ip := get_TCP_world_view_send_ips_matrix(ID, incoming_e_ID)

	fmt.Println("Server listening on ip: ", server_ip)
	ln, err := net.Listen("tcp", server_ip)
	if err != nil {
		fmt.Println("Error in TCP_server_setup", server_ip, err)
	}

	fmt.Println("Waiting for Accept:", server_ip)
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error in TCP_server_setup:", server_ip, err)
	}

	fmt.Println("Elevator ", incoming_e_ID, " connected to elevator ", ID, ". Setting ", incoming_e_ID, " to online")

	set_listen_dail_conn_matrix(ID, incoming_e_ID, conn)

	SetElevatorOnline(incoming_e_ID)

	set_trying_to_setup_matrix(ID, incoming_e_ID, false)

	defer ln.Close()
}

func TCP_client_setup(e_dailing_to_ID int) {
	set_trying_to_setup_matrix(e_dailing_to_ID, ID, true)

	client_ip := get_TCP_world_view_send_ips_matrix(e_dailing_to_ID, ID)

	for {
		fmt.Printf("Trying to dail to ip: %s\n", client_ip)
		conn, err := net.Dial("tcp", client_ip)
		if err != nil {
			fmt.Println("Dailing id ", e_dailing_to_ID, "failed, retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("Connected to", client_ip)

		set_listen_dail_conn_matrix(e_dailing_to_ID, ID, conn)

		SetElevatorOnline(e_dailing_to_ID) //setting status of connected elevator to online

		break
	}
	set_trying_to_setup_matrix(e_dailing_to_ID, ID, false)
}

func handle_receive(conn net.Conn, TCP_receive_channel chan masterSlave.Worldview_package, ID_of_connected_elevator int, i int, j int) {
	defer conn.Close()
	set_rescever_running_matrix(i, j, true)

	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(ID_of_connected_elevator))
	for {
		// Replace with actual receiving logic
		buffer := make([]byte, 1024)

		err := conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
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

		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving or timedout, closing receive goroutine and conn")
			SetElevatorOffline(ID_of_connected_elevator) //setting status of connected elevator to offline
			set_rescever_running_matrix(i, j, false)
			return
		}
		fmt.Printf("DATA MOTATT! ")

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
	send_world_view_package := masterSlave.New_Worldview_package(ID, fsm.GetElevatorStruct())
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

	//Finding package length
	packetLength := uint32(len(serialized_world_view_package)) //uint32 is 4 bytes

	// Sending to all online that are connected to us
	for connected_e_ID := ID + 1; connected_e_ID < Nr_of_elevators; connected_e_ID++ {
		if IsOnline(connected_e_ID) {

			//sending first packetLength, before actual packet. Preventing packet stacking
			server_conn := get_listen_dail_conn_matrix(ID, connected_e_ID)
			err = binary.Write(server_conn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost.")
				SetElevatorOffline(connected_e_ID) //setting status of connected elevator to offline
			}

			//writing acctual package
			_, err = server_conn.Write(serialized_world_view_package)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connected_e_ID) //setting status of connected elevator to offline
			}
		}
	}
	for connected_e_ID := 0; connected_e_ID < ID; connected_e_ID++ {
		if IsOnline(connected_e_ID) {
			//sending first packetLength, before actual packet. Preventing packet stacking
			client_conn := get_listen_dail_conn_matrix(connected_e_ID, ID)
			err = binary.Write(client_conn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost.")
				SetElevatorOffline(connected_e_ID) //setting status of connected elevator to offline
			}

			//writing acctual package
			_, err = client_conn.Write(serialized_world_view_package)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connected_e_ID) //setting status of connected elevator to offline
			}
		}
	}
}
