package connectivity

import "Driver-go/fsm"

var (
	cyclic_counter         = 0
	newButtonRequestMatrix [4][2]int
)

type WorldviewPackage struct {
	ElevatorID       int
	cyclic_counter   int
	Elevator         fsm.Elevator //refrence
	NewButtonRequest [4][2]int
}

func NewWorldviewPackage(ElevatorID int, elevator fsm.Elevator) WorldviewPackage {
	return WorldviewPackage{
		ElevatorID:       ElevatorID,
		cyclic_counter:   cyclic_counter,
		Elevator:         elevator,
		NewButtonRequest: newButtonRequestMatrix,
	}
}

func Increment_cyclic_counter() {
	cyclic_counter++
}

func Get_cyclic_counter() int {
	return cyclic_counter
}

func Set_cyclic_counter(c int) {
	cyclic_counter = c
}

func Get_newButtonRequestMatrix() [4][2]int {
	return newButtonRequestMatrix
}

func Set_newButtonRequestMatrix(floor int, button int) {
	newButtonRequestMatrix[floor][button] = 1
}
