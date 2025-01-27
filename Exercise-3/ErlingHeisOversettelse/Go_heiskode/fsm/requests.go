package fsm

import (
	"Driver-go/elevio"
	"fmt"
)

type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection
	behaviour ElevatorBehaviour
}



func requests_above(elevator Elevator) bool {
	for i := elevator.floor + 1; i < NumFloors; i++ {
		for j := 0; j < NumButtons ; j++ {
			if elevator.requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requests_below(elevator Elevator) bool {
	for i := 0; i < elevator.floor; i++ {
		for j := 0; j < NumButtons ; j++ {
			if elevator.requests[i][j] {
				return true
			}
		}
	}
	return false
}

func requests_here(elevator Elevator) bool {
	for j := 0; j < NumButtons ; j++ {
		if elevator.requests[elevator.floor][j] {
			return true
		}
	}
	return false
}


func requests_choose_direction(elevator Elevator) DirnBehaviourPair {
	switch elevator.dirn {
		
	}