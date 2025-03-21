package main

import (
	"Driver-go/connectivity"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"time"
)

const (
	Port_server_id0 = 15657
)

func main() {

	connect_to_elevatorserver()

	// Communication with elevator server setup
	drv_buttons, drv_floors, drv_obstr := elevio.Io_threds_setup()

	// Networking setup
	TCP_receive_channel, world_view_send_ticker, offline_update_chan := connectivity.Connectivity_setup()

	// Sets up timer
	timerTimeoutChan := fsm.Fsm_threds_setup()

	// Makes sure network connections have time to start properly
	time.Sleep(2000 * time.Millisecond)

	// Sets elevator to valid start possition
	fsm.SetElevatorToValidStartPossition()

	fmt.Println("Started!")

	// Stores the previous floor to detect floor changes
	prev_floor := -1

	// Logic loop for elevator and communication
	for {

		select {

		// Button press event
		case button_event := <-drv_buttons:
			fmt.Printf("\nButton event: %+v\n", button_event)

			// Starts order assignment if other elevators are online and it’s not a cab request
			if len(connectivity.Get_all_online_ids()) != 1 && button_event.Button != elevio.BT_Cab {
				connectivity.PrintIsOnline()
				connectivity.New_order(button_event)

			} else {

				// Handles request if no other elevators are online or it’s a cab request
				fmt.Println("No other online elevators or a cab call. Take order")
				fsm.Fsm_onRequestButtonPress(button_event.Floor, button_event.Button)
			}

		// Floor event
		case floor := <-drv_floors:
			fmt.Printf("Floor event: %+v\n", floor)

			// If elevator arrives at a different floor
			if floor != -1 && floor != prev_floor {
				fsm.Fsm_onFloorArrival(floor)
			}
			prev_floor = floor

		// Door timeout after 3 seconds
		case timer_bool := <-timerTimeoutChan:
			if timer_bool {
				fmt.Println("Door timeout")
				fsm.TimerStop()
				fsm.Fsm_onDoorTimeout()
			}

		// If there is an obstruction event
		case obstr_event_bool := <-drv_obstr:
			fmt.Println("Obstruction event toggle")
			fsm.SetObstructionStatus(obstr_event_bool)
			fsm.TimerStart(3)

		// World view ticker happens every 100 milliseconds
		case <-world_view_send_ticker:
			// Update lights and attempt to send world view
			connectivity.SetAllLights()
			connectivity.Send_world_view()

		// Incoming worldview package from another elevator
		case received_world_view := <-TCP_receive_channel:
			connectivity.Store_worldview(received_world_view.Elevator_ID, received_world_view)

			// If the received world view contains an order
			if received_world_view.Order_bool {
				fmt.Println("Order received")
				fsm.Fsm_onRequestButtonPress(received_world_view.Order.Floor, received_world_view.Order.Button)
			}

		// If an elevator goes offline, retrieve its ID and take over its orders
		case id_of_offline_elevator := <-offline_update_chan:
			fmt.Println("Elevator has disconnected. Running start backup")
			connectivity.Start_backup_process(id_of_offline_elevator)

		}

	}

}

func connect_to_elevatorserver() {
	// Setting up connection with elevator server

	var port int
	if connectivity.USE_IPS {
		//if USE_IPS true, use deafult port for elevator server
		port = Port_server_id0

	} else {
		// if USE_IPs false, use increasing port nr
		port = Port_server_id0 + connectivity.ID
	}
	ip := fmt.Sprintf("localhost:%d", port)
	fmt.Println("ID: ", connectivity.ID, ", ip: ", ip)
	elevio.Init(ip, fsm.NumFloors)
}
