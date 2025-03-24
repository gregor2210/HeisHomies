package main

import (
	"Driver-go/connectivity"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"time"
)

const (
	NUMFLOORS     = 4
	PortServerID0 = 15657
)

func main() {

	port := PortServerID0 + connectivity.ID
	ip := fmt.Sprintf("localhost:%d", port)
	fmt.Println("ID: ", connectivity.ID, ", ip: ", ip)
	elevio.Init(ip, NUMFLOORS)

	//var d elevio.MotorDirection = elevio.MotorUp
	//elevio.SetMotorDirection(d)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	//connectivity.TCP_setup()
	tcpReceiveChannel := make(chan connectivity.WorldviewPackage)
	//TCP_send_channel_listen := make(chan connectivity.WorldviewPackage)
	//TCP_send_channel_dail := make(chan connectivity.WorldviewPackage)
	go connectivity.TcpReceivingSetup(tcpReceiveChannel)

	// Go routine to send world view every second
	var worldViewSendTicker <-chan time.Time
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	worldViewSendTicker = ticker.C

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	//Order setup
	order_to_send_chan := make(chan connectivity.DoneProcessedOrder)
	connectivity.Order_setup(order_to_send_chan)
	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	//Online setup
	offlineUpdateChan := make(chan int)
	connectivity.OnlineSetup(offlineUpdateChan)
	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Started!")

	//inputPollRateMs := 25
	prevFloor := -1

	fsm.SetElevatorToValidStartPosition()

	for {
		select {
		// Kan enten få inn en ButtonEvent, en etasje (int) eller en obstruction
		case buttonEvent := <-drvbuttons: // Hvis det kommer en ButtonEvent {Floor, ButtonType} fra chanelen drvbuttons
			fmt.Println("Button event-------------------------------------------------------------------------")
			fmt.Printf("%+v\n", buttonEvent)

			if len(connectivity.GetAllOnlineIds()) != 1 && buttonEvent.Button != elevio.BtnCab {
				connectivity.PrintIsOnline()
				//This is start the prosses of finding the best elevator, only if there are other elevators online
				//This will also not run if it is a cab request
				priorityValue := fsm.CalculatePriorityValue(buttonEvent)
				connectivity.NewOrder(buttonEvent, priorityValue)
			} else {
				// If elevator do not see any other elevators are online. Do the request selfe
				fmt.Println("No other online elevators or a cab call. Take order")
				fsm.FsmOnRequestButtonPress(buttonEvent.Floor, buttonEvent.Button)
			}

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

			SendWorldviewPackage := connectivity.NewWorldviewPackage(connectivity.ID, fsm.GetElevatorStruct())
			connectivity.PrintOrderRequest(SendWorldviewPackage.Order_requeset)
			connectivity.PrintOrderRequest(SendWorldviewPackage.Order_response)

			//connectivity.PrintIsOnline()

			//case worldView := <-tcpReceiveChannel:
			//fmt.Println("World view reseved, PC:", worldView.ElevatorID, "\n")
			//fmt.Println("\n\n")
			//time.Sleep(500 * time.Duration(inputPollRateMs))

		case receivedWorldView := <-tcpReceiveChannel:
			//fmt.Println("World view reseved, PC:", receivedWorldView.ElevatorID)

			//storing worldview
			connectivity.StoreWorldview(receivedWorldView.ElevatorID, receivedWorldView)

			if receivedWorldView.OrderBool {
				fmt.Println("Order receved")
				fsm.FsmOnRequestButtonPress(receivedWorldView.Order.Floor, receivedWorldView.Order.Button)
			}

			//priorityValue := fsm.CalculatePriorityValue(receivedWorldView.)
			connectivity.Receved_order_requests(receivedWorldView.Order_requeset, receivedWorldView.ElevatorID) //Mulig vi kan flytte denne inn i conneciton pakka
			//fmt.Println("Receved_order_response, id: ", receivedWorldView.ElevatorID)
			connectivity.Receved_order_response(receivedWorldView.Order_response)

		case received_order := <-order_to_send_chan:
			if !connectivity.SendOrderToSpesificElevator(received_order) {
				fmt.Println("Failed to send order. taking it selfe")
				fsm.FsmOnRequestButtonPress(received_order.Order.Floor, received_order.Order.Button)
			}

		case idOfflineElevator := <-offlineUpdateChan:
			//When online staus of a elevator goes from online to offline. We get the id and start the backup prosess
			//THis will insure not lost calls
			connectivity.Start_backup_prosess(idOfflineElevator)

		}

	}

}
