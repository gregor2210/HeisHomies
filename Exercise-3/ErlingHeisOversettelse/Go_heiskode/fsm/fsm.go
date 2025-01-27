package fsm

import (
	"Driver-go/elevio"
	"fmt"
)

func Fsm_button_clicked_selecter(button elevio.ButtonType, floor int) {
	switch button {
	case elevio.BT_Cab:
		fsm_cab_button_pressed(floor)
	case elevio.BT_HallUp:
		fsm_hall_up_button_pressed(floor)
	default:
		fsm_hall_down_button_pressed(floor)
	}
}

func fsm_cab_button_pressed(floor int) {
	fmt.Println("Cab button pressed at floor ", floor)
}

func fsm_hall_up_button_pressed(floor int) {
	fmt.Println("Hall up button pressed at floor ", floor)
}

func fsm_hall_down_button_pressed(floor int) {
	fmt.Println("Hall down button pressed at floor ", floor)
}

var elevator Elevator

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

func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	fmt.Printf("\n\nfsm_onRequestButtonPress(%d, %s)\n", btn_floor)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			timer_start(elevator.doorOpenDuration_s)
		} else {
			elevator.requests[btn_floor][btn_type] = 1
		}
	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = 1
	case EB_Idle:
		elevator.requests[btn_floor][btn_type] = 1
		pair := requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			timer_start(elevator.config.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
		case EB_Idle:
		}
	}

	setAllLights(elevator)

	fmt.Println("\nNew state:")
}

func fsm_onFloorArrival(newFloor int) {
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

func fsmOnDoorTimeout() {
	elevator.floor = newFloor //usikker på newfloor
	switch ElevatorBehaviour {
	case EB_DoorOpen:
		//Usikker på hvordan DirnBehaviourPair er siden ikke er laget enda

		var pair DirnBehaviourPair = requestsChooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.dirn = pair.behaviour
		switch elevator.behaviour {
		case EB_DoorOpen:
			TimerStart(elevator.doorOpenDuration_s)
			//time.Sleep(time.Duration(elevator.config.doorOpenDuration_s) * time.Second))
			//Denne er ikke god nok, siden den låser mottak av ordre...
			elevator = requests_clearAtCurrentFloor(elevator)
			setAllLights(elevator)
			break
		case EB_Moving:
		case EB_Idle:
			outputDevice.doorLight(0)
			outputDevice.motorDirection(elevator.dirn)
			break
		}

	default:
		break
	}
}
