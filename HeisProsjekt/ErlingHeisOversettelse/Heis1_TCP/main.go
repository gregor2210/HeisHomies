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
	// Connect to elevator server
	connect_to_elevatorserver()
	// Communication with elevator server setup
	drv_buttons, drv_floors, drv_obstr := elevio.Io_threds_setup()

	// Nettworking settup
	TCP_receive_channel, world_view_send_ticker, offline_update_chan := connectivity.Connectivity_setup()

	// Sets up all fsm threds
	timerTimeoutChan := fsm.Fsm_threds_setup()

	// Makes sure network connections have time to start properly
	time.Sleep(2000 * time.Millisecond)

	// Sets elevator to valid start possition
	fsm.SetElevatorToValidStartPossition()

	fmt.Println("Started!")

	// Variabel for logic in forloop
	// Keeps track of the previuse floor
	prev_floor := -1

	// Logic loop for elevator and communication
	for {
		select {
		// Can either receive a ButtonEvent, a floor (int), or an obstruction
		case button_event := <-drv_buttons: // If a ButtonEvent {Floor, ButtonType} comes from the channel drv_buttons
			//fmt.Println("Button event-------------------------------------------------------------------------")
			fmt.Printf("\nButton event: %+v\n", button_event)

			if len(connectivity.Get_all_online_ids()) != 1 && button_event.Button != elevio.BT_Cab {
				connectivity.PrintIsOnline()
				// This starts the process of finding the best elevator, only if there are other elevators online
				// This will also not run if it is a cab request

				connectivity.New_order(button_event)
			} else {
				// If the elevator does not see any other elevators online, do the request itself
				fmt.Println("No other online elevators or a cab call. Take order")
				fsm.Fsm_onRequestButtonPress(button_event.Floor, button_event.Button)
			}

		case floor := <-drv_floors: // If a floor (int) comes from the channel drv_floors
			//fmt.Println("Floor event")
			fmt.Printf("Floor event: %+v\n", floor)
			if floor != -1 && floor != prev_floor {
				// If the elevator is at a floor and the floor is different from the previous floor
				fsm.Fsm_onFloorArrival(floor)
			}
			prev_floor = floor

		case timer_bool := <-timerTimeoutChan: // If a bool, True, comes from the channel timerTimeoutChan
			if timer_bool {
				fmt.Println("Door timeout")
				fsm.TimerStop()
				fsm.Fsm_onDoorTimeout()
			}

		case obstr_event_bool := <-drv_obstr: // If there is an obstruction event, BOOL
			fmt.Println("Obstruction event toggle")
			fsm.SetObstructionStatus(obstr_event_bool)
			fsm.TimerStart(3)

		case <-world_view_send_ticker: // World view ticker happens every x milliseconds
			// 1. Checks if lights should be turned off or on
			// 2. Attempts to send world view
			connectivity.SetAllLights()
			connectivity.Send_world_view()

		case received_world_view := <-TCP_receive_channel: // An incoming world view package, from other computers
			// Storing world view
			connectivity.Store_worldview(received_world_view.Elevator_ID, received_world_view)

			if received_world_view.Order_bool {
				fmt.Println("Order received")
				fsm.Fsm_onRequestButtonPress(received_world_view.Order.Floor, received_world_view.Order.Button)
			}

		case id_of_offline_elevator := <-offline_update_chan: // When an elevator goes from online to offline, this receives the now offline elevator id
			// When the online status of an elevator goes from online to offline, we get the id and start the backup process
			// This will ensure no lost calls
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
