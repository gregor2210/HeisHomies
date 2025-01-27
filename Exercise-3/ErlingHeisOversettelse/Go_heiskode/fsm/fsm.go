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
			timer_start(elevator.config.doorOpenDuration_s)
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
	elevator_print(elevator)
}
