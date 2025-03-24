package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"math"
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

func CalculatePriorityValue(buttonEvent elevio.ButtonEvent) int {
	requestFloor := buttonEvent.Floor
	//request_button := buttonEvent.Button

	//button point dir. minus is down
	//requst_button_point_dir := -1
	//if int(request_button) == 0 {
	//requst_button_point_dir = 1
	//}

	elevator := GetElevatorStruct()
	NumFloorsMinus1 := NumFloors - 1
	//Calculate how much this elevator wants this request.
	priorityValue := 2 * 10 * NumFloorsMinus1 // max value

	//DÅRLIG VERSJON, TAR ABSOLUTT AVSTAND
	deltaFloor := requestFloor - elevator.Floor
	//fmt.Println("---------------------------------------------------")
	//fmt.Println("request floor: ", requestFloor)
	//fmt.Println("elevator floor: ", elevator.Floor)
	//fmt.Println("Delta floor: ", deltaFloor)
	//fmt.Println("Priority value :", priorityValue)

	priorityValue = priorityValue - int(math.Abs(float64(deltaFloor)))*10

	//FUGERTE IKKE HELT
	/*
		//antall etasjer unna gir -10 poeng
		//deltaFloor gir minus value if requested floor is below
		deltaFloor := requestFloor - elevator.Floor

		//if elevator dosen ot have a moving dirn
		if int(elevator.Dirn) == 0 {
			priorityValue = priorityValue - int(math.Abs(float64(deltaFloor)))*10

			//HVis heis beveger seg mot request, og request peker i samme retning som heisens bevegelse
		} else if math.Copysign(1, float64(deltaFloor)) == math.Copysign(1, float64(elevator.Dirn)) {
			// deltaFloor have the same sign as elv direction. Meaning it is going towards the request)

			if math.Copysign(1, float64(requst_button_point_dir)) == math.Copysign(1, float64(elevator.Dirn)) {
				//Elevator moves towards request and in same direction as request
				priorityValue = priorityValue - int(math.Abs(float64(priorityValue)-float64(deltaFloor)))*10

			} else {
				//Elevator moves toward request, but request is in oposit riection
				if int(elevator.Dirn) < 0 {
					//elevator moves downward
					wortcase_down := elevator.Floor + requestFloor //goes to bottom, and turns direciton and moves up
					priorityValue = priorityValue - wortcase_down*10
				} else {
					//elevator moves up
					worstcase_up := (NumFloorsMinus1 - elevator.Floor) + (NumFloorsMinus1 - requestFloor)
					priorityValue = priorityValue - worstcase_up*10
				}
			}

		} else {
			// Elevator is not going towards the request

			if math.Copysign(1, float64(requst_button_point_dir)) == math.Copysign(1, float64(elevator.Dirn)) {
				//Elevator will not point in same drection after turn at a worstcase top

			} else {
				//Elevator moves toward request, but request is in oposit riection
				if int(elevator.Dirn) < 0 {
					//elevator moves downward
					nr_of_floors_traveled := elevator.Floor + NumFloorsMinus1 + requestFloor
					priorityValue = priorityValue - nr_of_floors_traveled*10
				} else {
					//elevator moves up
					nr_of_floors_traveled := (NumFloorsMinus1 - elevator.Floor) + NumFloorsMinus1 + requestFloor
					priorityValue = priorityValue - nr_of_floors_traveled*10

				}
			}

		}*/
	fmt.Println("Priority value:", priorityValue)
	//fmt.Println("---------------------------------------------------")
	return priorityValue
}
