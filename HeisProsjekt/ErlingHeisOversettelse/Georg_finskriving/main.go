package main

import (
	"Driver-go/connectivity"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"time"
)

const (
	NUMFLOORS       = 4
	Port_server_id0 = 15657
)

func main() {

	port := Port_server_id0 + connectivity.ID
	ip := fmt.Sprintf("localhost:%d", port)
	fmt.Println("ID: ", connectivity.ID, ", ip: ", ip)
	elevio.Init(ip, NUMFLOORS)
	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	//connectivity.TCP_setup()
	TCP_receive_channel := make(chan connectivity.Worldview_package)

	go connectivity.TCP_receving_setup(TCP_receive_channel)

	// Go routine to send world view every second
	var world_view_send_ticker <-chan time.Time
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	world_view_send_ticker = ticker.C

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	//Online setup
	offline_update_chan := make(chan int)
	connectivity.Online_setup(offline_update_chan)
	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Coms with server setup
	fsm.InitDriver()

	go fsm.PollTimerTimeout(timerTimeoutChan)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Started!")

	//inputPollRateMs := 25
	prev_floor := -1

	fsm.SetElevatorToValidStartPossition()

	for {
		select {
		// Kan enten få inn en ButtonEvent, en etasje (int) eller en obstruction
		case button_event := <-drv_buttons: // Hvis det kommer en ButtonEvent {Floor, ButtonType} fra chanelen drv_buttons
			fmt.Println("Button event-------------------------------------------------------------------------")
			fmt.Printf("%+v\n", button_event)

			if len(connectivity.Get_all_online_ids()) != 1 && button_event.Button != elevio.BT_Cab {
				connectivity.PrintIsOnline()
				//This is start the prosses of finding the best elevator, only if there are other elevators online
				//This will also not run if it is a cab request

				connectivity.New_order(button_event)
			} else {
				// If elevator do not see any other elevators are online. Do the request selfe
				fmt.Println("No other online elevators or a cab call. Take order")
				fsm.Fsm_onRequestButtonPress(button_event.Floor, button_event.Button)
			}

		case a := <-drv_floors: // Hvis det kommer en etasje (int) fra chanelen drv_floors
			fmt.Println("Floor event")
			fmt.Printf("%+v\n", a)
			if a != -1 && a != prev_floor { // Hvis heisen er i en etasje og etasjen er ulik den forrige etasjen
				fsm.Fsm_onFloorArrival(a)
			}
			prev_floor = a

		case a := <-timerTimeoutChan: // Hvis det kommer en bool, True, fra chanelen timerTimeoutChan
			if a {
				fmt.Println("Door timeout")
				fsm.TimerStop()
				fsm.Fsm_onDoorTimeout()
			}

		case a := <-drv_obstr:
			fmt.Println("Obstruction event toggle")
			fsm.SetObsructionStatus(a)
			fsm.TimerStart(3)

		case <-world_view_send_ticker:
			// 1. sjekker om lamper skal av eller på
			// 2. Prøver å sende worldview

			connectivity.SetAllLights()

			//fmt.Println("Sending world view")
			connectivity.Send_world_view()
			//connectivity.PrintIsOnline()

			//time.Sleep(500 * time.Duration(inputPollRateMs))

		case received_world_view := <-TCP_receive_channel:
			//fmt.Println("World view reseved, PC:", received_world_view.Elevator_ID)

			//storing worldview
			connectivity.Store_worldview(received_world_view.Elevator_ID, received_world_view)

			if received_world_view.Order_bool {
				fmt.Println("Order receved")
				fsm.Fsm_onRequestButtonPress(received_world_view.Order.Floor, received_world_view.Order.Button)
			}

		case id_of_offline_elevator := <-offline_update_chan:
			//When online staus of a elevator goes from online to offline. We get the id and start the backup prosess
			//THis will insure not lost calls
			fmt.Println("Running start backup")
			connectivity.Start_backup_prosess(id_of_offline_elevator)

		}

	}

}
