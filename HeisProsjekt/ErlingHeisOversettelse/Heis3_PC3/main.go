package main

import (
	"Driver-go/connectivity"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"time"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15659", numFloors)

	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

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

	// Go routine to send world view every second
	var world_view_send_ticker <-chan time.Time
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	world_view_send_ticker = ticker.C

	// Channel to receive world view
	world_view_resever_chan := make(chan fsm.Elevator) // WORLD VIEW IS ONLY A STRING TEMPORALRY
	go connectivity.Receive_elevator_world_view_distributor(world_view_resever_chan)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	fmt.Println("Started!")

	inputPollRateMs := 25
	prev_floor := -1

	fsm.SetElevatorToValidStartPossition()

	for {
		select {
		// Kan enten få inn en ButtonEvent, en etasje (int) eller en obstruction
		case a := <-drv_buttons: // Hvis det kommer en ButtonEvent {Floor, ButtonType} fra chanelen drv_buttons
			fmt.Println("Button event-------------------------------------------------------------------------")
			fmt.Printf("%+v\n", a)
			fsm.Fsm_onRequestButtonPress(a.Floor, a.Button)

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

		//fmt.Printf("%+v\n", a)

		//case a := <-drv_stop:
		//fmt.Printf("%+v\n", a)

		case world_view := <-world_view_resever_chan:
			fmt.Println("World view reseved:")
			fsm.PrintElevator(world_view)

		case <-world_view_send_ticker:
			connectivity.Send_elevator_world_view()
		}

		time.Sleep(500 * time.Duration(inputPollRateMs))

	}

}
