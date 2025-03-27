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

	// Sets elevator to valid start possition
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

	wg.Wait()

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

func elevatorFunctionality(drvFloors <-chan int, motorErrorChan <-chan bool, timerTimeOutChan <-chan bool,
	drvObstr <-chan bool, tcpReceiveChannel <-chan connectivity.WorldviewPackage, obstrErrorChan <-chan bool) {
	// Loop for the critical elevator functionality that need to be handled without delay

	// Stores the previous floor to detect floor changes
	prevFloor := -1

	// Logic loop for elevator and communication
	// The main loop that make sure only one event is handled at a time
	for {

		select {

		// Floor event
		case floor := <-drvFloors:
			fmt.Printf("Floor event: %+v\n", floor)

			// If elevator arrives at a different floor
			if floor != -1 && floor != prevFloor {
				fsm.FsmOnFloorArrival(floor)
			}
			prevFloor = floor

			// if an elevaotr get to a floor. the motor works!
			connectivity.SetSelfOnline()
			fsm.SetElevatorMotorError(false)

		// Motor Error detected
		case errorBool := <-motorErrorChan:
			// There is an error with the motor
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

		case errorBool := <-obstrErrorChan:
			// Obstruction has been on for to long
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
	// Loop for the critical network functionality that can use multible seconds to execute
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
