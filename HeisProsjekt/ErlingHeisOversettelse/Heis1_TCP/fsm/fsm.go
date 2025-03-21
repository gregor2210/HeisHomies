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

	switch elevator.Behaviour {
	case EB_DoorOpen:
		// Button pressed at current floor
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			TimerStart(elevator.DoorOpenDuration_s) // Restart door timer

		} else {
			elevator.Requests[btn_floor][btn_type] = true // // Add request to queue
		}
	case EB_Moving:
		elevator.Requests[btn_floor][btn_type] = true // In motion, so only add to the queue

	case EB_Idle:
		elevator.Requests[btn_floor][btn_type] = true
		var pair DirnBehaviourPair = requests_chooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			TimerStart(elevator.DoorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)

		case EB_Moving:
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn))

		case EB_Idle:
		}
	}

	setAllLights(elevator)

}

func Fsm_onFloorArrival(newFloor int) {
	elevator_mutex.Lock()
	defer elevator_mutex.Unlock()
	elevator.Floor = newFloor
	elevio.SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:

		// Stop if request in direction, cab call, or no more requests
		if requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			TimerStart(elevator.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
	}

}

func Fsm_onDoorTimeout() {
	elevator_mutex.Lock()
	defer elevator_mutex.Unlock()

	switch elevator.Behaviour {
	case EB_DoorOpen:

		var pair DirnBehaviourPair = requests_chooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			// Restart door timer
			TimerStart(elevator.DoorOpenDuration_s)

			elevator = requests_clearAtCurrentFloor(elevator)

			setAllLights(elevator)

		case EB_Moving:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn))

		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
		}

	default:
		// Nothing to do if the state is not EB_DoorOpen
	}

}
