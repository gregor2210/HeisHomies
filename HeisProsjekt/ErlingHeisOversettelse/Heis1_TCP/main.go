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
	PortServerID0 = 15657
)

func main() {

	connectToElevatorserver()

	// Communication with elevator server setup
	drvButtons, drvFloors, drvObstr := elevio.InitIOHandling()

	// Networking setup
	tcpReceiveChannel, worldViewSendTicker, offlineUpdateChan := connectivity.ConnectivitySetup()

	// Sets up timer
	timerTimeOutChan, motorErrorChan, obstrErrorChan := fsm.FsmThreadsSetup()

	// Makes sure network connections have time to start properly
	time.Sleep(2000 * time.Millisecond)

	// Sets elevator to valid start position
	fsm.SetElevatorToValidStartPosition()

	fmt.Println("Started!")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		elevatorFunctionality(drvFloors, motorErrorChan, timerTimeOutChan, drvObstr, tcpReceiveChannel, obstrErrorChan)
	}()

	go func() {
		defer wg.Done()
		networkFunctionality(drvButtons, worldViewSendTicker, offlineUpdateChan)
	}()

	// Make sure the elevator start running, set a call on current floor

	wg.Wait()

}

// Setting up connection with elevator server
func connectToElevatorserver() {

	var port int

	//if UseIPs true, use deafult port for elevator server
	if connectivity.UseIPs {

		port = PortServerID0

	} else {
		// if UseIPs false, use increasing port number
		port = PortServerID0 + connectivity.ID
	}
	ip := fmt.Sprintf("localhost:%d", port)
	fmt.Println("ID: ", connectivity.ID, ", ip: ", ip)
	elevio.Init(ip, fsm.NumFloors)
}

// Loop for the critical elevator functionality that needs to be handled without delay
func elevatorFunctionality(drvFloors <-chan int, motorErrorChan <-chan bool, timerTimeOutChan <-chan bool,
	drvObstr <-chan bool, tcpReceiveChannel <-chan connectivity.WorldviewPackage, obstrErrorChan <-chan bool) {

	// Stores the previous floor to detect floor changes
	prevFloor := -1

	// Main loop for elevator and network logic
	// Ensures only one event is handled at a time
	for {

		select {

		// Floor event
		case floor := <-drvFloors:
			fmt.Printf("Floor event: %+v\n", floor)

			// If elevator arrives at a different floor
			if floor != -1 && floor != prevFloor {
				fsm.FsmOnFloorArrival(floor)
			}

			// -1 is initial condition: simulate call to start movement
			if prevFloor == -1 {
				fmt.Println("Starting elevator movement")
				fsm.FsmOnRequestButtonPress(floor, 2)
			}
			prevFloor = floor

			// Floor reached, confirm motor is working
			connectivity.SetSelfOnline()
			fsm.SetElevatorMotorError(false)

		// Motor Error detected
		case errorBool := <-motorErrorChan:
			fmt.Println("Motor error Error")
			// if errorBool == True and not the online elevator
			if errorBool {
				fmt.Println("Elevator has motor problems. Running start backup")
				fsm.SetElevatorMotorError(true)
				connectivity.StartErrorProcess()

				// Resets motor direction of elevator
				elevio.SetMotorDirection(elevio.MotorStop)
				if prevFloor <= 1 {
					elevio.SetMotorDirection(elevio.MotorUp)
				} else {
					elevio.SetMotorDirection(elevio.MotorDown)
				}
			}
		// Obstruction active for too long
		case errorBool := <-obstrErrorChan:

			fmt.Println("Obstruction error Error")
			// if errorBool == True and not the online elevator
			if errorBool && connectivity.SelfOnlyOnline() { // If the elevator is alone
				fmt.Println("Elevator has obstruction problems but is alone. Starting new timer")
				fsm.StartObstrTimer()
				connectivity.SetSelfOffline()
				connectivity.CloseAllConnections()

			} else if errorBool && !connectivity.SelfOnlyOnline() { // If the elevator is not alone
				fmt.Println("Elevator has obstruction problems. Setting selfe offline and closing all connections")
				connectivity.StartErrorProcess()

			}

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
			motorError := fsm.GetElevatorMotorError()
			if !motorError { // If there is no motor error
				connectivity.SetSelfOnline()
			}
			fsm.TimerDoorStart(3)

		// Incoming worldview package from another elevator
		case receivedWorldView := <-tcpReceiveChannel:
			connectivity.StoreWorldview(receivedWorldView.ElevatorID, receivedWorldView)

			// If the received world view contains an order
			if receivedWorldView.OrderBool {
				fmt.Println("Order received")
				fsm.FsmOnRequestButtonPress(receivedWorldView.Order.Floor, receivedWorldView.Order.Button)
			}
		}
		connectivity.SetAllLights()
	}
}

func networkFunctionality(drvButtons <-chan elevio.ButtonEvent, worldViewSendTicker <-chan time.Time, offlineUpdateChan <-chan int) {
	// Loop for the critical network functionality that can use multiple seconds to execute
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
				fmt.Println("No other online elevators or a cab call. Taking order")
				fsm.FsmOnRequestButtonPress(buttonEvent.Floor, buttonEvent.Button)
			}

		// World view ticker happens every 100 milliseconds
		case <-worldViewSendTicker:
			// Attempt to send world view
			connectivity.SendWorldView()

		// If an elevator goes offline, retrieve its ID and take over its orders
		case idOfflineElevator := <-offlineUpdateChan:
			fmt.Println("Elevator has disconnected. Running start backup")
			connectivity.StartBackupProcess(idOfflineElevator)
		}
		connectivity.SetAllLights()
	}
}
