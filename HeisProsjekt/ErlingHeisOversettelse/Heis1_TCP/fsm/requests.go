package fsm

import (
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}

func requests_above(elevator Elevator) bool {
	for i := elevator.Floor + 1; i < NumFloors; i++ {
		for j := 0; j < NumButtons; j++ {
			if elevator.Requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requests_below(elevator Elevator) bool {
	for i := 0; i < elevator.Floor; i++ {
		for j := 0; j < NumButtons; j++ {
			if elevator.Requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requests_here(elevator Elevator) bool {
	for j := 0; j < NumButtons; j++ {
		if elevator.Requests[elevator.Floor][j] {
			return true
		}
	}
	return false
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
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
	switch e.Dirn {
	case D_Down:
		return e.Requests[e.Floor][elevio.BT_HallDown] ||
			e.Requests[e.Floor][elevio.BT_Cab] ||
			!requests_below(e)
	case D_Up:
		return e.Requests[e.Floor][elevio.BT_HallUp] ||
			e.Requests[e.Floor][elevio.BT_Cab] ||
			!requests_above(e)
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

// Returns true if elevator is at a floor with a pressed button
func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	return e.Floor == btn_floor &&
		((e.Dirn == D_Up && btn_type == elevio.BT_HallUp) ||
			(e.Dirn == D_Down && btn_type == elevio.BT_HallDown) ||
			e.Dirn == D_Stop ||
			btn_type == elevio.BT_Cab)
}

// Clears requests at current floor based on elevator direction
func requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.Requests[e.Floor][elevio.BT_Cab] = false
	switch e.Dirn {
	case D_Up:
		if !requests_above(e) && !e.Requests[e.Floor][elevio.BT_HallUp] {
			e.Requests[e.Floor][elevio.BT_HallDown] = false
		}
		e.Requests[e.Floor][elevio.BT_HallUp] = false
	case D_Down:
		if !requests_below(e) && !e.Requests[e.Floor][elevio.BT_HallDown] {
			e.Requests[e.Floor][elevio.BT_HallUp] = false
		}
		e.Requests[e.Floor][elevio.BT_HallDown] = false

	case D_Stop:
		fallthrough

	default:
		e.Requests[e.Floor][elevio.BT_HallUp] = false
		e.Requests[e.Floor][elevio.BT_HallDown] = false
	}
	return e
}
