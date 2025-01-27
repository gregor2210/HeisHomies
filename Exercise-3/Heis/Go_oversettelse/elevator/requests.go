package elevator

import (
	"project/elevator"
)

// Sjekker om det finnes forespørsler over den nåværende etasjen
func requestsAbove(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < elevator.N_FLOORS; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

// Sjekker om det finnes forespørsler under den nåværende etasjen
func requestsBelow(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

// Sjekker om det finnes forespørsler i den nåværende etasjen
func requestsHere(e elevator.Elevator) bool {
	for btn := 0; btn < elevator.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] == 1 {
			return true
		}
	}
	return false
}

// Velger retning for heisen basert på forespørsler og nåværende retning
func requestsChooseDirection(e elevator.Elevator) elevator.DirnBehaviourPair {
	switch e.Dirn {

	case elevator.D_Up:
		if requestsAbove(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Up,
				Behaviour: elevator.EB_Moving,
			}
		} else if requestsHere(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Down,
				Behaviour: elevator.EB_DoorOpen,
			}
		} else if requestsBelow(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Down,
				Behaviour: elevator.EB_Moving,
			}
		} else {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Stop,
				Behaviour: elevator.EB_Idle,
			}
		}

	case elevator.D_Down:
		if requestsBelow(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Down,
				Behaviour: elevator.EB_Moving,
			}
		} else if requestsHere(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Up,
				Behaviour: elevator.EB_DoorOpen,
			}
		} else if requestsAbove(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Up,
				Behaviour: elevator.EB_Moving,
			}
		} else {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Stop,
				Behaviour: elevator.EB_Idle,
			}
		}

	case elevator.D_Stop:
		if requestsHere(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Stop,
				Behaviour: elevator.EB_DoorOpen,
			}
		} else if requestsAbove(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Up,
				Behaviour: elevator.EB_Moving,
			}
		} else if requestsBelow(e) {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Down,
				Behaviour: elevator.EB_Moving,
			}
		} else {
			return elevator.DirnBehaviourPair{
				Dirn:      elevator.D_Stop,
				Behaviour: elevator.EB_Idle,
			}
		}

	default:
		return elevator.DirnBehaviourPair{
			Dirn:      elevator.D_Stop,
			Behaviour: elevator.EB_Idle,
		}
	}
}

// Sjekker om heisen bør stoppe i den nåværende etasjen
func requestsShouldStop(e elevator.Elevator) bool {
	switch e.Dirn {

	case elevator.D_Down:
		return e.Requests[e.Floor][elevator.B_HallDown] ||
			e.Requests[e.Floor][elevator.B_Cab] ||
			!requestsBelow(e)

	case elevator.D_Up:
		return e.Requests[e.Floor][elevator.B_HallUp] ||
			e.Requests[e.Floor][elevator.B_Cab] ||
			!requestsAbove(e)

	case elevator.D_Stop:
		return true

	default:
		return false
	}
}

// Sjekker om en forespørsel bør fjernes umiddelbart
func requestsShouldClearImmediately(e elevator.Elevator, btnFloor int, btnType elevator.Button) bool {
	switch e.Config.ClearRequestVariant {

	case elevator.CV_All:
		return e.Floor == btnFloor

	case elevator.CV_InDirn:
		return e.Floor == btnFloor &&
			((e.Dirn == elevator.D_Up && btnType == elevator.B_HallUp) ||
				(e.Dirn == elevator.D_Down && btnType == elevator.B_HallDown) ||
				e.Dirn == elevator.D_Stop || btnType == elevator.B_Cab)

	default:
		return false
	}
}

// Fjerner forespørsler i den nåværende etasjen basert på konfigurasjonen
func requestsClearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	switch e.Config.ClearRequestVariant {

	case elevator.CV_All:
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = 0
		}

	case elevator.CV_InDirn:
		e.Requests[e.Floor][elevator.B_Cab] = 0

		switch e.Dirn {
		case elevator.D_Up:
			if !requestsAbove(e) && e.Requests[e.Floor][elevator.B_HallUp] == 0 {
				e.Requests[e.Floor][elevator.B_HallDown] = 0
			}
			e.Requests[e.Floor][elevator.B_HallUp] = 0

		case elevator.D_Down:
			if !requestsBelow(e) && e.Requests[e.Floor][elevator.B_HallDown] == 0 {
				e.Requests[e.Floor][elevator.B_HallUp] = 0
			}
			e.Requests[e.Floor][elevator.B_HallDown] = 0

		case elevator.D_Stop:
			e.Requests[e.Floor][elevator.B_HallUp] = 0
			e.Requests[e.Floor][elevator.B_HallDown] = 0
		}
	}
	return e
}