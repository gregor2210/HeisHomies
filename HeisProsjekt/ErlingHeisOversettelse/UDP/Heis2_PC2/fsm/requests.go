package fsm

import (
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}

func requestsAbove(elevator Elevator) bool {
	for i := elevator.Floor + 1; i < NumFloors; i++ {
		for j := 0; j < NumButtons; j++ {
			if elevator.Requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(elevator Elevator) bool {
	for i := 0; i < elevator.Floor; i++ {
		for j := 0; j < NumButtons; j++ {
			if elevator.Requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requestsHere(elevator Elevator) bool {
	for j := 0; j < NumButtons; j++ {
		if elevator.Requests[elevator.Floor][j] {
			return true
		}
	}
	return false
}

func requestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case DirUp:
		if requestsAbove(e) {
			return DirnBehaviourPair{DirUp, ElevMoving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{DirStop, ElevDoorOpen}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{DirDown, ElevMoving}
		} else {
			return DirnBehaviourPair{DirStop, ElevIdle}
		}
	case DirDown:
		if requestsBelow(e) {
			return DirnBehaviourPair{DirDown, ElevMoving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{DirStop, ElevDoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{DirUp, ElevMoving}
		} else {
			return DirnBehaviourPair{DirStop, ElevIdle}
		}
	case DirStop:
		if requestsHere(e) {
			return DirnBehaviourPair{DirStop, ElevDoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{DirUp, ElevMoving}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{DirDown, ElevMoving}
		} else {
			return DirnBehaviourPair{DirStop, ElevIdle}
		}
	default:
		return DirnBehaviourPair{DirStop, ElevIdle}
	}
}

func requestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case DirDown:
		return e.Requests[e.Floor][elevio.BtnHallDown] ||
			e.Requests[e.Floor][elevio.BtnCab] ||
			!requestsBelow(e)
	case DirUp:
		return e.Requests[e.Floor][elevio.BtnHallUp] ||
			e.Requests[e.Floor][elevio.BtnCab] ||
			!requestsAbove(e)
	case DirStop:
		fallthrough
	default:
		return true
	}
}

// Denne funksjonen sjekker om heisen er i en etasje og om det er en knapp som er trykket inn i den etasjen. Dersom det er det, vil den returnere true.
func requestsShouldClearImmediately(e Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return e.Floor == btnFloor &&
		((e.Dirn == DirUp && btnType == elevio.BtnHallUp) ||
			(e.Dirn == DirDown && btnType == elevio.BtnHallDown) ||
			e.Dirn == DirStop ||
			btnType == elevio.BtnCab)
}

func requestsClearAtCurrentFloor(e Elevator) Elevator {
	e.Requests[e.Floor][elevio.BtnCab] = false
	switch e.Dirn {
	case DirUp:
		if !requestsAbove(e) && !e.Requests[e.Floor][elevio.BtnHallUp] {
			e.Requests[e.Floor][elevio.BtnHallDown] = false
		}
		e.Requests[e.Floor][elevio.BtnHallUp] = false
	case DirDown:
		if !requestsBelow(e) && !e.Requests[e.Floor][elevio.BtnHallDown] {
			e.Requests[e.Floor][elevio.BtnHallUp] = false
		}
		e.Requests[e.Floor][elevio.BtnHallDown] = false

	case DirStop:
		fallthrough

	default:
		e.Requests[e.Floor][elevio.BtnHallUp] = false
		e.Requests[e.Floor][elevio.BtnHallDown] = false
	}
	return e
}
