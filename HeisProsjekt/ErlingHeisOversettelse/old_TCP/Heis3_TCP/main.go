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

	//var d elevio.MotorDirection = elevio.MotorUp
	//elevio.SetMotorDirection(d)

	drvbuttons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drv_stop := make(chan bool)

	// Channel to receive timer TimeOut events
	timerTimeOutChan := make(chan bool)

	go elevio.PollButtons(drvbuttons)
	//søker igjennom alle etasjene og sjekker alle typer knapper for den etasjen.
	//Den sjekker ved å sende en tpc getbutton(etasje, knappetype) og får tilbake true/false. Dersom dette er anderledes enn fra forigje gang den sjekket og den nå er nå true.
	//Skriver den til chanelen drvbuttons. Den skriver da en ButtonEvent (struct) med etasje og knappetype.
	go elevio.PollFloorSensor(drvFloors)
	//Sjekker om heisen er i en etasje og den etasjen er ulik det den var sist gang den sjekket. Dersom den er det, skriver den til chanelen drvFloors. Den skriver da etasjen heisen er i, i form av en int.
	go elevio.PollObstructionSwitch(drvObstr)
	//Sjekker om det er en obstruction i heisen. Dersom statuesn på obstruction endrer seg så skriver den true til chanelen drvObstr. Dersom det ikke er obstruction, skriver den false.
	go elevio.PollStopButton(drv_stop)
	//Sjekker om stopknappen er trykket inn. Den vil skrive true eller false til chanelen drv_stop når statusen endrer seg.

	go fsm.PollTimerTimeOut(timerTimeOutChan)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	//connectivity.TCP_setup()
	tcpReceiveChannel := make(chan connectivity.WorldviewPackage)
	TCP_send_channel_listen := make(chan connectivity.WorldviewPackage)
	TCP_send_channel_dail := make(chan connectivity.WorldviewPackage)
	go connectivity.TcpReceivingSetup(tcpReceiveChannel, TCP_send_channel_listen, TCP_send_channel_dail)

	// Go routine to send world view every second
	var worldViewSendTicker <-chan time.Time
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	worldViewSendTicker = ticker.C

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Started!")

	inputPollRateMs := 25
	prevFloor := -1

	fsm.SetElevatorToValidStartPosition()

	for {
		select {
		// Kan enten få inn en ButtonEvent, en etasje (int) eller en obstruction
		case a := <-drvbuttons: // Hvis det kommer en ButtonEvent {Floor, ButtonType} fra chanelen drvbuttons
			fmt.Println("Button event-------------------------------------------------------------------------")
			fmt.Printf("%+v\n", a)
			fsm.FsmOnRequestButtonPress(a.Floor, a.Button)

		case a := <-drvFloors: // Hvis det kommer en etasje (int) fra chanelen drvFloors
			fmt.Println("Floor event")
			fmt.Printf("%+v\n", a)
			if a != -1 && a != prevFloor { // Hvis heisen er i en etasje og etasjen er ulik den forrige etasjen
				fsm.FsmOnFloorArrival(a)
			}
			prevFloor = a

		case a := <-timerTimeOutChan: // Hvis det kommer en bool, True, fra chanelen timerTimeOutChan
			if a {
				fmt.Println("Door TimeOut")
				fsm.TimerStop()
				fsm.FsmOnDoorTimeOut()
			}

		case a := <-drvObstr:
			fmt.Println("Obstruction event toggle")
			fsm.SetObsructionStatus(a)
			fsm.TimerStart(3)

		case <-worldViewSendTicker:
			//fmt.Println("Sending world view")
			connectivity.SendWorldView()
			//connectivity.PrintIsOnline()

			//case worldView := <-tcpReceiveChannel:
			//fmt.Println("World view reseved, PC:", worldView.ElevatorID, "\n")
			//fmt.Println("\n\n")
			time.Sleep(500 * time.Duration(inputPollRateMs))

		case recived_worldView := <-tcpReceiveChannel:
			fmt.Println("World view reseved, PC:", recived_worldView.ElevatorID)
		}

	}

}
