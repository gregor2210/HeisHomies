package main

import (
	"Driver-go/connectivity"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"sync"
	"time"
)

const (
	PortServerID0    = 15657
	DelayBeforeStart = 2000 // Milliseconds
)

func main() {

	connectToElevatorserver()

	// Elevator IO setup 
	drvButtons, drvFloors, drvObstr := elevio.ElevatorIoSetup()

	// Networking setup
	tcpReceiveChannel, worldViewSendTicker, offlineUpdateChan := connectivity.ConnectivitySetup()

	// Timer setup
	timerTimeOutChan, motorErrorChan, obstrErrorChan := fsm.FsmThreadsSetup()

	// Makes sure network connections have time to start properly
	time.Sleep(DelayBeforeStart * time.Millisecond)

	// Sets elevator to valid start position
	elevio.SetElevatorToValidStartPosition()

	fmt.Println("Started!")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		criticalElevatorFunctionality(drvFloors, motorErrorChan, timerTimeOutChan, drvObstr, tcpReceiveChannel, obstrErrorChan)
	}()

	go func() {
		defer wg.Done()
		networkHandeling(drvButtons, worldViewSendTicker, offlineUpdateChan)
	}()

	wg.Wait()

}

// Setting up connection with elevator server
func connectToElevatorserver() {

	var port int

	if connectivity.UseIPs {
		port = PortServerID0

	} else {
		port = PortServerID0 + connectivity.ID
	}
	ip := fmt.Sprintf("localhost:%d", port)
	fmt.Println("ID: ", connectivity.ID, ", ip: ", ip)
	elevio.Init(ip, fsm.NumFloors)
}

// Loop for the critical elevator functionality that needs to be handled without delay
func criticalElevatorFunctionality(drvFloors <-chan int, motorErrorChan <-chan bool, timerTimeOutChan <-chan bool,
	drvObstr <-chan bool, tcpReceiveChannel <-chan connectivity.WorldviewPackage, obstrErrorChan <-chan bool) {


	prevFloor := -1 // Initial condition: no floor registered yet

	// Main loop: handles elevator and network logic, one event at a time
	for {

		select {

		// Floor event
		case floor := <-drvFloors:
			fmt.Printf("Floor event: %+v\n", floor)

			// If elevator arrives at a different floor
			if floor != -1 && floor != prevFloor {
				fsm.FsmOnFloorArrival(floor)
			}

			// On startup: if no previous floor, simulate a button press to trigger movement
			if prevFloor == -1 {
				fmt.Println("Starting elevator movement")
				fsm.FsmOnRequestButtonPress(floor, 2)
			}
			prevFloor = floor

			// Floor reached, confirm motor is working
			connectivity.SetSelfOnline()
			fsm.SetElevatorMotorError(false)

		// Incoming worldview package from another elevator
		case receivedWorldView := <-tcpReceiveChannel:
			connectivity.StoreWorldview(receivedWorldView.ElevatorID, receivedWorldView)
			if receivedWorldView.OrderBool {
				fmt.Println("Order received")
				fsm.FsmOnRequestButtonPress(receivedWorldView.Order.Floor, receivedWorldView.Order.Button)
			}
		// Motor Error detected
		case errorBool := <-motorErrorChan:
			fmt.Println("Motor error Error")
			if errorBool {
				fmt.Println("Elevator has motor problems. Running start backup")
				fsm.SetElevatorMotorError(true)
				connectivity.DisableComunicaton()
				elevio.MoveElevatorToFloor(prevFloor)
			}
		// Obstruction active for too long
		case errorBool := <-obstrErrorChan:
			if errorBool {
				fmt.Println("Elevator has obstruction problems Disable communication")
				fsm.StartObstrTimer()
				connectivity.DisableComunicaton()
			}

		// Obstruction event detected
		case obstrEventBool := <-drvObstr:
			fmt.Println("Obstruction event toggle")
			fsm.SetObstructionStatus(obstrEventBool)

			motorError := fsm.GetElevatorMotorError()
			if !motorError { 
				connectivity.SetSelfOnline()
			}
			fsm.TimerDoorStart(3)

		// Door TimeOut after 3 seconds
		case timerBool := <-timerTimeOutChan:
			if timerBool {
				fmt.Println("Door TimeOut")
				fsm.TimerStop()
				fsm.FsmOnDoorTimeOut()
			}
		}
		connectivity.SetAllLights()
	}
}

// Loop for the critical network functionality that can use multiple seconds to execute
func networkHandeling(drvButtons <-chan elevio.ButtonEvent, worldViewSendTicker <-chan time.Time, offlineUpdateChan <-chan int) {

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
				fmt.Println("No other online elevators or a cab call. Taking order")
				fsm.FsmOnRequestButtonPress(buttonEvent.Floor, buttonEvent.Button)
			}

		// World view ticker happens every 100 milliseconds
		case <-worldViewSendTicker:
			connectivity.SendWorldView()

		// If an elevator goes offline, retrieve its ID and take over its orders
		case idOfflineElevator := <-offlineUpdateChan:
			fmt.Println("Elevator has disconnected. Running start backup")
			connectivity.StartBackupProcess(idOfflineElevator)
		}
		connectivity.SetAllLights()
	}
}
