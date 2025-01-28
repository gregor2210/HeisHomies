package fsm

import (
	"Driver-go/elevio"
	"fmt"
)

var elevator Elevator = NewElevator()

//ElevatorOUtputDevice er den utdelte go driverern!

func setAllLights(elevator Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, elevator.requests[floor][btn])
		}
	}
}

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	fmt.Printf("\n\nfsm_onRequestButtonPress(%d)\n", btn_floor)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			TimerStart(elevator.doorOpenDuration_s)
		} else {
			elevator.requests[btn_floor][btn_type] = true
		}
	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = true

	case EB_Idle:
		elevator.requests[btn_floor][btn_type] = true
		var pair DirnBehaviourPair = requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			TimerStart(elevator.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)

		case EB_Moving:
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.dirn))

		case EB_Idle:
		}
	}

	setAllLights(elevator)

	fmt.Println("\nNew state:")
}

func Fsm_onFloorArrival(newFloor int) {
	fmt.Printf("\n\nfsm_onFloorArrival(%d)\n", newFloor)
	elevator.floor = newFloor

	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			TimerStart(elevator.doorOpenDuration_s)
			setAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}
	default:
	}

	fmt.Println("\nNew state:")
}

func Fsm_onDoorTimeout() {
	fmt.Printf("\n\nfsm_onDoorTimeout()\n")

	switch elevator.behaviour {
	case EB_DoorOpen:
		var pair DirnBehaviourPair = requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			TimerStart(elevator.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)
			setAllLights(elevator)
		case EB_Moving:
		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.dirn))
		}

	default:
	}

	fmt.Println("\nNew state:")
}
