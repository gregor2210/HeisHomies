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

	// setting up connection with elevator server
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
	elevio.Init(ip, NUMFLOORS)

	// Setting up TCP connection loop.
	TCP_receive_channel := make(chan connectivity.Worldview_package)

	go connectivity.TCP_receving_setup(TCP_receive_channel)

	// Go routine to send world view every second
	var world_view_send_ticker <-chan time.Time
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	world_view_send_ticker = ticker.C

	//Online setup
	offline_update_chan := make(chan int)
	connectivity.Online_setup(offline_update_chan)

	// Communication with elevator server setup
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// Channel to receive timer timeout events
	timerTimeoutChan := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	//søker igjennom alle etasjene og sjekker alle typer knapper for den etasjen.
	//Den sjekker ved å sende en tpc getbutton(etasje, knappetype) og får tilbake true/false. Dersom dette er anderledes enn fra forigje gang den sjekket og den nå er nå true.
	//Skriver den til chanelen drv_buttons. Den skriver da en ButtonEvent (struct) med etasje og knappetype.
	go elevio.PollFloorSensor(drv_floors)
	//Sjekker om heisen er i en etasje og den etasjen er ulik det den var sist gang den sjekket. Dersom den er det, skriver den til chanelen drv_floors. Den skriver da etasjen heisen er i, i form av en int.
	go elevio.PollObstructionSwitch(drv_obstr)
	//Sjekker om det er en obstruction i heisen. Dersom statuesn på obstruction endrer seg så skriver den true til chanelen drv_obstr. Dersom det ikke er obstruction, skriver den false.
	go elevio.PollStopButton(drv_stop)
	//Sjekker om stopknappen er trykket inn. Den vil skrive true eller false til chanelen drv_stop når statusen endrer seg.
	go fsm.PollTimerTimeout(timerTimeoutChan)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	time.Sleep(2000 * time.Millisecond)

	// Sets elevator to valid start possition
	fsm.SetElevatorToValidStartPossition()

	fmt.Println("Started!")

	// Variabel for logic in forloop
	// Keeps track of the previuse floor
	prev_floor := -1

	// Logic loop for elevator and comunication
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
			if a != -1 && a != prev_floor {
				// Hvis heisen er i en etasje og etasjen er ulik den forrige etasjen
				fsm.Fsm_onFloorArrival(a)
			}
			prev_floor = a

		case a := <-timerTimeoutChan: // Hvis det kommer en bool, True, fra chanelen timerTimeoutChan
			if a {
				fmt.Println("Door timeout")
				fsm.TimerStop()
				fsm.Fsm_onDoorTimeout()
			}

		case a := <-drv_obstr: // If there is an obstuction event, BOOL
			fmt.Println("Obstruction event toggle")
			fsm.SetObsructionStatus(a)
			fsm.TimerStart(3)

		case <-world_view_send_ticker: // World view ticker happens every x milliseconds
			// 1. sjekker om lamper skal av eller på
			// 2. Prøver å sende worldview
			connectivity.SetAllLights()
			connectivity.Send_world_view()

		case received_world_view := <-TCP_receive_channel: // A incoming worldview package, from other computers
			//storing worldview
			connectivity.Store_worldview(received_world_view.Elevator_ID, received_world_view)

			if received_world_view.Order_bool {
				fmt.Println("Order receved")
				fsm.Fsm_onRequestButtonPress(received_world_view.Order.Floor, received_world_view.Order.Button)
			}

		case id_of_offline_elevator := <-offline_update_chan: // When an elevator goes from online to ofline, this recives the now ofline elevator id
			//When online staus of a elevator goes from online to offline. We get the id and start the backup prosess
			//This will insure not lost calls
			fmt.Println("Running start backup")
			connectivity.Start_backup_prosess(id_of_offline_elevator)

		}

	}

}
