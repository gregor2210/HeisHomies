package fsm

import (
	"Driver-go/elevio"
	"sync"
)

var (
	elevator       Elevator = NewElevator()
	elevator_mutex sync.Mutex
)

func GetElevatorStruct() Elevator {
	elevator_mutex.Lock()
	defer elevator_mutex.Unlock()
	return elevator
}

//ElevatorOUtputDevice er den utdelte go driverern!

func setAllLights(elevator Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, elevator.Requests[floor][btn])
		}
	}
}

// func fsm_onInitBetweenFloors() {
// 	elevio.SetMotorDirection(elevio.MD_Down)
// 	elevator.dirn = elevio.MD_Down
// 	elevator.behaviour = EB_Moving
// }

func Fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	elevator_mutex.Lock()
	defer elevator_mutex.Unlock()
	//fmt.Printf("\n\nfsm_onRequestButtonPress(%d)\n", btn_floor)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) { // If the elevator is already at the floor and the button is pressed
			TimerStart(elevator.DoorOpenDuration_s) // Restart the door timer
		} else {
			elevator.Requests[btn_floor][btn_type] = true // Otherwise, add the request to the queue
		}
	case EB_Moving:
		elevator.Requests[btn_floor][btn_type] = true // In motion, so only add to the queue

	case EB_Idle:
		elevator.Requests[btn_floor][btn_type] = true                   // The elevator is idle, must determine direction and start moving
		var pair DirnBehaviourPair = requests_chooseDirection(elevator) // Choose direction based on requests
		elevator.Dirn = pair.dirn                                       // Update direction
		elevator.Behaviour = pair.behaviour                             // Update state
		switch pair.behaviour {
		case EB_DoorOpen: // If the elevator should stop at the current floor
			elevio.SetDoorOpenLamp(true)
			TimerStart(elevator.DoorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)

		case EB_Moving: // If the elevator should start moving
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn)) // Start the motor

		case EB_Idle:
		}
	}

	setAllLights(elevator) // Update the light indicators

	//fmt.Println("\nNew state:")
}

func Fsm_onFloorArrival(newFloor int) {
	elevator_mutex.Lock()
	defer elevator_mutex.Unlock()
	//fmt.Printf("\n\nfsm_onFloorArrival(%d)\n", newFloor)
	elevator.Floor = newFloor

	elevio.SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) { // If the elevator should stop at the current floor, either due to a request in the correct direction, a cab call, or no more requests
			elevio.SetMotorDirection(elevio.MD_Stop)          // Stop the motor
			elevio.SetDoorOpenLamp(true)                      // Open the door
			elevator = requests_clearAtCurrentFloor(elevator) // Clear requests at the current floor
			TimerStart(elevator.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
	}

	//fmt.Println("\nNew state:")
}

func Fsm_onDoorTimeout() {
	elevator_mutex.Lock()
	defer elevator_mutex.Unlock()
	//fmt.Printf("fsm_onDoorTimeout()")

	switch elevator.Behaviour {
	case EB_DoorOpen:
		// Choose direction based on requests
		var pair DirnBehaviourPair = requests_chooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			// Restart the door timer
			TimerStart(elevator.DoorOpenDuration_s)

			// Clear requests at the current floor
			elevator = requests_clearAtCurrentFloor(elevator)

			// Update all light indicators
			setAllLights(elevator)

		case EB_Moving:
			// Start the motor
			elevio.SetDoorOpenLamp(false) // Turn off the door indicator
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn))
			//fmt.Printf("Motor started moving in direction: %v\n", elevator.Dirn)

		case EB_Idle:
			// No more requests, set to idle
			elevio.SetDoorOpenLamp(false)
		}

	default:
		// Nothing to do if the state is not EB_DoorOpen
	}

	//fmt.Println("\nNew state:")
}
