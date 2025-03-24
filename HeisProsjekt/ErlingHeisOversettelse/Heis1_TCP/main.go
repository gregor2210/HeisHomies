package main

import (
	"Driver-go/connectivity"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"time"
)

const (
	PortServerID0 = 15657
)

func main() {

	connectToElevatorserver()

	// Communication with elevator server setup
	drvButtons, drvFloors, drvObstr := elevio.InitIOHandling()

	// Networking setup
	tcpReceiveChannel, worldViewSendTicker, offlineUpdateChan := connectivity.ConnectivitySetup()

	// Sets up timer
	timerTimeOutChan := fsm.FsmThreadsSetup()

	// Makes sure network connections have time to start properly
	time.Sleep(2000 * time.Millisecond)

	// Sets elevator to valid start possition
	fsm.SetElevatorToValidStartPosition()

	fmt.Println("Started!")

	// Stores the previous floor to detect floor changes
	prevFloor := -1

	// Logic loop for elevator and communication
	for {

		select {

		// Button press event
		case buttonEvent := <-drvButtons:
			fmt.Printf("\nButton event: %+v\n", buttonEvent)

			// Starts order assignment if other elevators are online and it’s not a cab request
			if len(connectivity.GetAllOnlineIds()) != 1 && buttonEvent.Button != elevio.BtnCab {
				connectivity.PrintIsOnline()
				connectivity.NewOrder(buttonEvent)

			} else {

				// Handles request if no other elevators are online or it’s a cab request
				fmt.Println("No other online elevators or a cab call. Take order")
				fsm.FsmOnRequestButtonPress(buttonEvent.Floor, buttonEvent.Button)
			}

		// Floor event
		case floor := <-drvFloors:
			fmt.Printf("Floor event: %+v\n", floor)

			// If elevator arrives at a different floor
			if floor != -1 && floor != prevFloor {
				fsm.FsmOnFloorArrival(floor)
			}
			prevFloor = floor

		// Door TimeOut after 3 seconds
		case timerBool := <-timerTimeOutChan:
			if timerBool {
				fmt.Println("Door TimeOut")
				fsm.TimerStop()
				fsm.FsmOnDoorTimeOut()
			}

		// If there is an obstruction event
		case obstrEventBool := <-drvObstr:
			fmt.Println("Obstruction event toggle")
			fsm.SetObstructionStatus(obstrEventBool)
			fsm.TimerStart(3)

		// World view ticker happens every 100 milliseconds
		case <-worldViewSendTicker:
			// Update lights and attempt to send world view
			connectivity.SetAllLights()
			connectivity.SendWorldView()

		// Incoming worldview package from another elevator
		case receivedWorldView := <-tcpReceiveChannel:
			connectivity.StoreWorldview(receivedWorldView.ElevatorID, receivedWorldView)

			// If the received world view contains an order
			if receivedWorldView.OrderBool {
				fmt.Println("Order received")
				fsm.FsmOnRequestButtonPress(receivedWorldView.Order.Floor, receivedWorldView.Order.Button)
			}

		// If an elevator goes offline, retrieve its ID and take over its orders
		case idOfflineElevator := <-offlineUpdateChan:
			fmt.Println("Elevator has disconnected. Running start backup")
			connectivity.StartBackupProcess(idOfflineElevator)

		}

	}

}

func connectToElevatorserver() {
	// Setting up connection with elevator server

	var port int
	if connectivity.UseIPs {
		//if UseIPs true, use deafult port for elevator server
		port = PortServerID0

	} else {
		// if UseIPs false, use increasing port nr
		port = PortServerID0 + connectivity.ID
	}
	ip := fmt.Sprintf("localhost:%d", port)
	fmt.Println("ID: ", connectivity.ID, ", ip: ", ip)
	elevio.Init(ip, fsm.NumFloors)
}
