package connectivity

import "Driver-go/fsm"

type WorldviewPackage struct {
	ElevatorID int
	Elevator   fsm.Elevator //refrence
}

func NewWorldviewPackage(ElevatorID int, elevator fsm.Elevator) WorldviewPackage {
	return WorldviewPackage{
		ElevatorID: ElevatorID,
		Elevator:   elevator,
	}
}
