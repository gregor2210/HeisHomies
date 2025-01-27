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
