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
	// Matrix of TCP IPs where rows are listeners and columns are dialers
	TCP_world_view_send_ips_matrix [NR_OF_ELEVATORS - 1][NR_OF_ELEVATORS]string

	// Matrix of active TCP connections between elevators
	listen_dail_conn_matrix [NR_OF_ELEVATORS - 1][NR_OF_ELEVATORS]net.Conn

	// Tracks active connection attempts between elevators
	trying_to_setup_matrix = [NR_OF_ELEVATORS - 1][NR_OF_ELEVATORS]bool{}

	// Matrix indicating active receiver goroutines
	receiver_running_matrix = [NR_OF_ELEVATORS - 1][NR_OF_ELEVATORS]bool{}

	// Mutex for the matrixes
	mu_world_view_send_ips_matrix sync.Mutex
	mu_listen_dail_conn_matrix    sync.Mutex
	mu_trying_to_setup_matrix     sync.Mutex
	mu_receiver_running_matrix    sync.Mutex
)

func init() { // runs when imported

	// This function setup the TCP_world_view_send_ips_matrix
	// If USE_IPS is true, use network IPs

	if USE_IPS {
		if NR_OF_ELEVATORS > len(IPs) {
			log.Fatal("NR_OF_ELEVATORS larger then amount of IPs")
		}

		for i := 0; i < NR_OF_ELEVATORS-1; i++ {
			for j := i + 1; j < NR_OF_ELEVATORS; j++ {
				ip := IPs[i] + ":80" + fmt.Sprint(i) + fmt.Sprint(j)
				set_TCP_world_view_send_ips_matrix(i, j, ip)
			}
		}

	} else {

		// If USE_IPS is false, set IPs to localhost for local testing
		for i := 0; i < NR_OF_ELEVATORS-1; i++ {
			for j := i + 1; j < NR_OF_ELEVATORS; j++ {
				ip := "localhost:80" + fmt.Sprint(i) + fmt.Sprint(j)
				set_TCP_world_view_send_ips_matrix(i, j, ip)
			}
		}
	}
}

// Thread-safe access functions for the connection matrices
func get_TCP_world_view_send_ips_matrix(i int, j int) string {
	mu_world_view_send_ips_matrix.Lock()
	defer mu_world_view_send_ips_matrix.Unlock()
	return TCP_world_view_send_ips_matrix[i][j]
}

func set_TCP_world_view_send_ips_matrix(i int, j int, ip string) {
	mu_world_view_send_ips_matrix.Lock()
	defer mu_world_view_send_ips_matrix.Unlock()
	TCP_world_view_send_ips_matrix[i][j] = ip
}

func get_listen_dail_conn_matrix(i int, j int) net.Conn {
	mu_listen_dail_conn_matrix.Lock()
	defer mu_listen_dail_conn_matrix.Unlock()
	return listen_dail_conn_matrix[i][j]
}

func set_listen_dail_conn_matrix(i int, j int, conn net.Conn) {
	mu_listen_dail_conn_matrix.Lock()
	defer mu_listen_dail_conn_matrix.Unlock()
	listen_dail_conn_matrix[i][j] = conn
}

func get_trying_to_setup_matrix(i int, j int) bool {
	mu_trying_to_setup_matrix.Lock()
	defer mu_trying_to_setup_matrix.Unlock()
	return trying_to_setup_matrix[i][j]
}

func set_trying_to_setup_matrix(i int, j int, b bool) {
	mu_trying_to_setup_matrix.Lock()
	defer mu_trying_to_setup_matrix.Unlock()
	trying_to_setup_matrix[i][j] = b
}

func get_receiver_running_matrix(i int, j int) bool {
	mu_receiver_running_matrix.Lock()
	defer mu_receiver_running_matrix.Unlock()
	return receiver_running_matrix[i][j]
}

func set_receiver_running_matrix(i int, j int, b bool) {
	mu_receiver_running_matrix.Lock()
	defer mu_receiver_running_matrix.Unlock()
	receiver_running_matrix[i][j] = b
}

// Serialize the struct
func serialize_elevator(wv Worldview_package) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wv)
	return buf.Bytes(), err
}

// Deserialize the incomming bytes
func deserialize_elevator(data []byte) (Worldview_package, error) {
	var wv Worldview_package
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wv)
	return wv, err
}

func TCP_receving_setup(TCP_receive_channel chan Worldview_package) {
	// Every 2 seconds, attempts to:
	// 1. Set up TCP servers and clients
	// 2. Establish missing connections
	// 3. Start receiver goroutines for active connections

	loop_timer := 2

	fmt.Println("Starting TCP receving setup")
	for {

		// Server setup
		for j := ID + 1; j < NR_OF_ELEVATORS; j++ {
			if !get_trying_to_setup_matrix(ID, j) && !IsOnline(j) && !get_receiver_running_matrix(ID, j) {
				fmt.Println("Starting up tcp_server_setup. ", ID, "listening for: ", j)
				go tcp_server_setup(j)
			}

			// Server receiver setup
			if IsOnline(j) && !get_receiver_running_matrix(ID, j) {
				fmt.Println("Starting handle_receive for connected elevator")
				go handle_receive(get_listen_dail_conn_matrix(ID, j), TCP_receive_channel, j, ID, j)
			}
		}

		//Client setup
		for i := 0; i < ID; i++ {
			//Client_conn_setup
			if !get_trying_to_setup_matrix(i, ID) && !IsOnline(i) && !get_receiver_running_matrix(i, ID) {
				fmt.Println("Starting up tcp_client_setup. ", ID, "dialing to: ", i)
				go tcp_client_setup(i)
			}

			//Client rescever setup
			if IsOnline(i) && !get_receiver_running_matrix(i, ID) {
				fmt.Println("Starting handle_receive for elevator we are connected to")
				go handle_receive(get_listen_dail_conn_matrix(i, ID), TCP_receive_channel, i, i, ID)
			}

		}
		time.Sleep(time.Duration(loop_timer) * time.Second)

	}

}

// Setting up server for self (ID) to listen to elevator (incoming_e_ID)
func tcp_server_setup(incoming_e_ID int) {

	set_trying_to_setup_matrix(ID, incoming_e_ID, true)

	server_ip := get_TCP_world_view_send_ips_matrix(ID, incoming_e_ID)

	fmt.Println("Server listening on ip: ", server_ip)
	ln, err := net.Listen("tcp", server_ip)
	if err != nil {
		fmt.Println("Error in tcp_server_setup", server_ip, err)
	}

	fmt.Println("Waiting for Accept:", server_ip)
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error in tcp_server_setup:", server_ip, err)
	}

	// Set no delay to true

	fmt.Println("Elevator ", incoming_e_ID, " connected to elevator ", ID, ". Setting ", incoming_e_ID, " to online")

	set_listen_dail_conn_matrix(ID, incoming_e_ID, conn)

	SetElevatorOnline(incoming_e_ID)

	set_trying_to_setup_matrix(ID, incoming_e_ID, false)

	defer ln.Close()
}

// Setting up client for self (ID) to dial a server (e_dailing_to_ID)
func tcp_client_setup(e_dailing_to_ID int) {

	set_trying_to_setup_matrix(e_dailing_to_ID, ID, true)

	client_ip := get_TCP_world_view_send_ips_matrix(e_dailing_to_ID, ID)

	for {

		// Will try to dial every second until it connects
		fmt.Printf("Trying to dail to ip: %s\n", client_ip)
		conn, err := net.Dial("tcp", client_ip)
		if err != nil {
			fmt.Println("Dailing id ", e_dailing_to_ID, "failed, retrying in 2 seconds...")
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Println("Connected to", client_ip)

		set_listen_dail_conn_matrix(e_dailing_to_ID, ID, conn)

		SetElevatorOnline(e_dailing_to_ID)

		break
	}
	set_trying_to_setup_matrix(e_dailing_to_ID, ID, false)
}

func handle_receive(conn net.Conn, TCP_receive_channel chan Worldview_package, ID_of_connected_elevator int, i int, j int) {
	// Runs as a goroutine to handle incoming messages:
	// 1. Sets a read deadline. If no data is received within TIMEOUT, the connection is closed.
	// 2. Reads the length of the incoming packet.
	// 3. Reads the packet of the specified length.
	// 4. Deserializes and sends the worldview package to TCP_receive_channel.

	defer conn.Close()
	set_receiver_running_matrix(i, j, true)

	fmt.Println("HANDLE RECEIVE STARTED, ID: " + fmt.Sprint(ID_of_connected_elevator))
	for {
		// Replace with actual receiving logic

		// Setting read deadline
		err := conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
		if err != nil {
			fmt.Println("Conn not open")
			SetElevatorOffline(ID_of_connected_elevator)
			set_receiver_running_matrix(i, j, false)
			return
		}

		// Read packet length
		var packetLength uint32
		err = binary.Read(conn, binary.BigEndian, &packetLength)
		if err != nil {
			SetElevatorOffline(ID_of_connected_elevator)
			fmt.Println("Failed to read packetLength:", err)
			set_receiver_running_matrix(i, j, false)
			return
		}

		// Read incomming worldview packet (bytes)
		buffer := make([]byte, packetLength)
		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving or timedout, closing receive goroutine and conn")
			SetElevatorOffline(ID_of_connected_elevator)
			set_receiver_running_matrix(i, j, false)
			return
		}

		// Deserialize the buffer to worldview package
		receved_world_view_package, err := deserialize_elevator(buffer)
		if err != nil {
			log.Fatal("failed to deserialize:", err)
		}

		// Store backup worldview from incomming elevator
		Store_worldview(receved_world_view_package.Elevator_ID, receved_world_view_package)

		TCP_receive_channel <- receved_world_view_package
	}

}

func Send_world_view() {
	// Tries to send a worldview package to all online elevators
	// First to those connected to us (servers), then to those we are connected to (clients)
	send_world_view_package := New_Worldview_package(ID, fsm.GetElevatorStruct())
	serialized_world_view_package, err := serialize_elevator(send_world_view_package)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}

	if len(serialized_world_view_package) > PACKAGE_SIZE {
		log.Fatal("error: serialized data too large")
	}

	//Finding package length
	packetLength := uint32(len(serialized_world_view_package)) //uint32 is 4 bytes

	// Send to elevators where we are server
	for connected_e_ID := ID + 1; connected_e_ID < NR_OF_ELEVATORS; connected_e_ID++ {
		if IsOnline(connected_e_ID) {

			// Send packet length first to avoid message stacking
			server_conn := get_listen_dail_conn_matrix(ID, connected_e_ID)

			// Start write session with timeout
			err = server_conn.SetWriteDeadline(time.Now().Add(TIMEOUT * time.Second))
			if err != nil {
				fmt.Println("Failed to set write deadline for server write:", err, connected_e_ID)
				SetElevatorOffline(connected_e_ID)
				server_conn.Close()
				continue

			}
			// Write packet length
			err = binary.Write(server_conn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost.")
				SetElevatorOffline(connected_e_ID) //setting status of connected elevator to offline
				server_conn.Close()
				continue
			}

			// Write actual package
			_, err = server_conn.Write(serialized_world_view_package)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connected_e_ID) //setting status of connected elevator to offline
				fmt.Println("Elevator was set to ofline!!!!! 123")
				server_conn.Close()
				continue
			}

			// Disable write deadline after transmission
			server_conn.SetWriteDeadline(time.Time{})
		}
	}

	// Send to elevators where we are client
	for connected_e_ID := 0; connected_e_ID < ID; connected_e_ID++ {
		if IsOnline(connected_e_ID) {

			client_conn := get_listen_dail_conn_matrix(connected_e_ID, ID)

			err = client_conn.SetWriteDeadline(time.Now().Add(TIMEOUT * time.Second))
			if err != nil {
				fmt.Println("Failed to set write deadline for client write:", err, connected_e_ID)
				SetElevatorOffline(connected_e_ID)
				client_conn.Close()
				continue
			}

			err = binary.Write(client_conn, binary.BigEndian, packetLength)
			if err != nil {
				fmt.Println("Error sending packetlength to connected elevator, connection lost.")
				SetElevatorOffline(connected_e_ID)
				client_conn.Close()
				continue
			}

			_, err = client_conn.Write(serialized_world_view_package)
			if err != nil {
				fmt.Println("Error sending, connection lost.")
				SetElevatorOffline(connected_e_ID)
				client_conn.Close()
				continue
			}

			client_conn.SetWriteDeadline(time.Time{})
		}
	}
}

func Send_order_to_spesific_elevator(recever_e int, order elevio.ButtonEvent) bool {
	// Find the correct connection, regardless of whether this elevator is dialing or listening
	var conn net.Conn
	if IsOnline(recever_e) {
		if ID < (NR_OF_ELEVATORS-1) && get_listen_dail_conn_matrix(ID, recever_e) != nil {
			conn = get_listen_dail_conn_matrix(ID, recever_e)

		} else if recever_e < (NR_OF_ELEVATORS-1) && get_listen_dail_conn_matrix(recever_e, ID) != nil {
			conn = get_listen_dail_conn_matrix(recever_e, ID)
		} else {
			fmt.Println("No valid conn to send ORDER")
			return false
		}
	} else {
		return false
	}

	send_world_view_package := New_Worldview_package(ID, fsm.GetElevatorStruct())
	send_world_view_package.Order_bool = true
	send_world_view_package.Order = order

	serialized_world_view_package, err := serialize_elevator(send_world_view_package)
	if err != nil {
		log.Fatal("failed to serialize:", err)
	}

	// Set the write deadline for both write operations (2 seconds)
	err = conn.SetWriteDeadline(time.Now().Add(TIMEOUT * time.Second))
	if err != nil {
		fmt.Println("Failed to set write deadline:", err)
		return false
	}

	// Finding package length
	packetLength := uint32(len(serialized_world_view_package)) // uint32 is 4 bytes

	err = binary.Write(conn, binary.BigEndian, packetLength)
	if err != nil {
		fmt.Println("Error sending packetlength for ORDER to connected elevator, connection lost or timedout.")
		SetElevatorOffline(recever_e)
		conn.Close()
		return false
	}
	fmt.Println("Success sending ORDER packetlength")

	// Writing actual package
	_, err = conn.Write(serialized_world_view_package)
	if err != nil {
		fmt.Println("Error sending ORDER, connection lost.  or timedout")
		SetElevatorOffline(recever_e)
		conn.Close()
		return false
	}

	// Disable SetwriteDeadline
	conn.SetWriteDeadline(time.Time{})

	//fmt.Println("Success sending ORDER")

	//everyting worked!
	return true
}
