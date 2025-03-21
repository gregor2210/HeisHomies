package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

const NumFloors int = 4
const NumButtons int = 3

func GetMotorDirectionFromDirn(dirn Dirn) elevio.MotorDirection {
	switch dirn {
	case D_Up:
		return elevio.MD_Up
	case D_Down:
		return elevio.MD_Down
	case D_Stop:
		return elevio.MD_Stop
	default:
		return elevio.MD_Stop
	}
}

type Dirn int

const (
	D_Up   Dirn = 1
	D_Down Dirn = -1
	D_Stop Dirn = 0
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2
)

type Elevator struct {
	ID                 int
	Floor              int
	Dirn               Dirn
	Behaviour          ElevatorBehaviour
	Requests           [NumFloors][NumButtons]bool
	DoorOpenDuration_s float64
	Obstruction        bool
}

// Elevator initializer function
func NewElevator() Elevator {
	var elevator_setup Elevator = Elevator{
		ID:                 0,
		Floor:              -1,                            // Uninitialized floor
		Dirn:               D_Stop,                        // Not moving
		Behaviour:          EB_Idle,                       // Idle state
		Requests:           [NumFloors][NumButtons]bool{}, // No requests initially
		DoorOpenDuration_s: 3.0,
		Obstruction:        false,
	}

	return elevator_setup
}

func PrintElevator(elevator Elevator) {
	fmt.Printf("\n\nElevator:\n")
	fmt.Printf("ID: %d\n", elevator.ID)
	fmt.Printf("Floor: %d\n", elevator.Floor)
	fmt.Printf("Direction: %v\n", elevator.Dirn)
	fmt.Printf("Behaviour: %v\n", elevator.Behaviour)
	fmt.Printf("Requests: %v\n", elevator.Requests)
	fmt.Printf("Obstruction: %v\n", elevator.Obstruction)
}

func SetObstructionStatus(status bool) {
	elevator.Obstruction = status
}

func SetElevatorToValidStartPossition() {
	fmt.Println("Elevator initialized")
	for {
		if elevio.GetFloor() == -1 {
			elevio.SetMotorDirection(elevio.MD_Down)
		} else {
			elevio.SetMotorDirection(elevio.MD_Stop)
			break
		}
		time.Sleep(_pollRate)

	}
}
