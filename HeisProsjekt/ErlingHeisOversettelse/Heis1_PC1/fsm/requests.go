package fsm

import (
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}

func requests_above(elevator Elevator) bool {
	for i := elevator.floor + 1; i < NumFloors; i++ {
		for j := 0; j < NumButtons; j++ {
			if elevator.requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requests_below(elevator Elevator) bool {
	for i := 0; i < elevator.floor; i++ {
		for j := 0; j < NumButtons; j++ {
			if elevator.requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requests_here(elevator Elevator) bool {
	for j := 0; j < NumButtons; j++ {
		if elevator.requests[elevator.floor][j] {
			return true
		}
	}
	return false
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requests_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if requests_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Stop:
		if requests_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}

func requests_shouldStop(e Elevator) bool {
	switch e.dirn {
	case D_Down:
		return e.requests[e.floor][elevio.BT_HallDown] ||
			e.requests[e.floor][elevio.BT_Cab] ||
			!requests_below(e)
	case D_Up:
		return e.requests[e.floor][elevio.BT_HallUp] ||
			e.requests[e.floor][elevio.BT_Cab] ||
			!requests_above(e)
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

// Denne funksjonen sjekker om heisen er i en etasje og om det er en knapp som er trykket inn i den etasjen. Dersom det er det, vil den returnere true.
func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	return e.floor == btn_floor &&
		((e.dirn == D_Up && btn_type == elevio.BT_HallUp) ||
			(e.dirn == D_Down && btn_type == elevio.BT_HallDown) ||
			e.dirn == D_Stop ||
			btn_type == elevio.BT_Cab)
}

func requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.requests[e.floor][elevio.BT_Cab] = false
	switch e.dirn {
	case D_Up:
		if !requests_above(e) && !e.requests[e.floor][elevio.BT_HallUp] {
			e.requests[e.floor][elevio.BT_HallDown] = false
		}
		e.requests[e.floor][elevio.BT_HallUp] = false
	case D_Down:
		if !requests_below(e) && !e.requests[e.floor][elevio.BT_HallDown] {
			e.requests[e.floor][elevio.BT_HallUp] = false
		}
		e.requests[e.floor][elevio.BT_HallDown] = false

	case D_Stop:
		fallthrough

	default:
		e.requests[e.floor][elevio.BT_HallUp] = false
		e.requests[e.floor][elevio.BT_HallDown] = false
	}
	return e
}
