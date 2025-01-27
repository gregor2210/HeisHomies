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
//Chat
func fsmOnDoorTimeout() {
	switch elevator.behaviour {
	case EB_DoorOpen:
		pair := requestsChooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			timerStart(elevator.config.doorOpenDuration_s)
			elevator = requestsClearAtCurrentFloor(elevator)
			setAllLights(elevator)
		case EB_Moving, EB_Idle:
			outputDevice.doorLight(0)
			outputDevice.motorDirection(elevator.dirn)
		}
	}
}


//Min
func fsmOnDoorTimeout(){
	elevator.floor=newFloor//usikker på newfloor
	switch(ElevatorBehaviour){
	case EB_DoorOpen
	//Usikker på hvordan DirnBehaviourPair er siden ikke er laget enda

	DirnBehaviourPair pair = requestsChooseDirection(elevator)
	elevator.dirn = pair.dirn
	elevator.dirn = pair.behaviour
		switch(elevator.behaviour){
		case EB_DoorOpen:
			TimerStart(elevator.config.doorOpenDuration_s)
			//time.Sleep(time.Duration(elevator.config.doorOpenDuration_s) * time.Second))
            //Denne er ikke god nok, siden den låser mottak av ordre...
			elevator = requests_clearAtCurrentFloor(elevator)
            setAllLights(elevator)
			break;
		case EB_Moving
		case EB_Idle
			outputDevice.doorLight(0)
			outputDevice.motorDirection(elevator.dirn)
			break;
		}

	default:
		break;
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
