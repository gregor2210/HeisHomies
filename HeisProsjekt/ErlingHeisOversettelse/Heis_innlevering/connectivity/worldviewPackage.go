package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
)

type WorldviewPackage struct {
	ElevatorID int
	Elevator   fsm.Elevator //refrence
	Order      elevio.ButtonEvent
	OrderBool  bool
}

func NewWorldviewPackage(ElevatorID int, elevator_ fsm.Elevator) WorldviewPackage {
	return WorldviewPackage{
		ElevatorID: ElevatorID,
		Elevator:   elevator_,
		OrderBool:  false,
	}
}

func PrintWorldview(worldView WorldviewPackage) {
	fmt.Print("WorldView")
	fmt.Println("Elevator ID: ", worldView.ElevatorID)

}
