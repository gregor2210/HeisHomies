package elevio

import (
	"fmt"
	"time"
)

func MoveElevatorToFloor(prevFloor int) {
	// Resets motor direction of elevator
	SetMotorDirection(MotorStop)
	if prevFloor <= 1 {
		SetMotorDirection(MotorUp)
	} else {
		SetMotorDirection(MotorDown)
	}
}

func SetElevatorToValidStartPosition() {
	fmt.Println("Elevator initialized")
	for {
		if GetFloor() == -1 {
			SetMotorDirection(MotorDown)
		} else {
			SetMotorDirection(MotorStop)
			break
		}
		time.Sleep(_pollRate)

	}
}
