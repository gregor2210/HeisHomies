package fsm

import (
	"Driver-go/elevio"
	"sync"
)

var (
	elevator      Elevator = NewElevator()
	elevatorMutex sync.Mutex
)

func GetElevatorStruct() Elevator {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()
	return elevator
}

func SetElevatorMotorError(motorError bool) {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()
	elevator.MotorError = motorError
}

func GetElevatorMotorError() bool {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()
	return elevator.MotorError
}

func setElevatorObtruction(obstruction bool) {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()
	elevator.Obstruction = obstruction
}

func setAllLights(elevator Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		elevio.SetButtonLamp(elevio.ButtonType(2), floor, elevator.Requests[floor][2])
	}
}

func FsmOnRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()

	switch elevator.Behaviour {
	case ElevDoorOpen:
		// Button pressed at current floor
		if requestsShouldClearImmediately(elevator, btnFloor, btnType) {
			TimerDoorStart(elevator.DoorOpenDuration_s) // Restart door timer

		} else {
			elevator.Requests[btnFloor][btnType] = true // // Add request to queue
		}
	case ElevMoving:
		elevator.Requests[btnFloor][btnType] = true // In motion, so only add to the queue

	case ElevIdle:
		elevator.Requests[btnFloor][btnType] = true
		var pair DirnBehaviourPair = requestsChooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch pair.behaviour {
		case ElevDoorOpen:
			elevio.SetDoorOpenLamp(true)
			TimerDoorStart(elevator.DoorOpenDuration_s)
			elevator = requestsClearAtCurrentFloor(elevator)

		case ElevMoving:
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn))
			StartMotorErrorTimer(elevator)

		case ElevIdle:
		}
	}
	storeCabRequests(elevator)
	setAllLights(elevator)

}

func FsmOnFloorArrival(newFloor int) {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()
	elevator.Floor = newFloor
	elevio.SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case ElevMoving:

		// Stop if request in direction, cab call, or no more requests
		if requestsShouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MotorStop)
			StopMotorErrorTimer()
			elevio.SetDoorOpenLamp(true)
			elevator = requestsClearAtCurrentFloor(elevator)
			TimerDoorStart(elevator.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = ElevDoorOpen
		}
	default:
	}
	storeCabRequests(elevator)

}

func FsmOnDoorTimeOut() {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()

	switch elevator.Behaviour {
	case ElevDoorOpen:

		var pair DirnBehaviourPair = requestsChooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch elevator.Behaviour {
		case ElevDoorOpen:
			// Restart door timer
			TimerDoorStart(elevator.DoorOpenDuration_s)

			elevator = requestsClearAtCurrentFloor(elevator)

			setAllLights(elevator)

		case ElevMoving:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn))
			StartMotorErrorTimer(elevator)

		case ElevIdle:
			elevio.SetDoorOpenLamp(false)
		}

	default:
		// Nothing to do if the state is not ElevDoorOpen
	}

}
