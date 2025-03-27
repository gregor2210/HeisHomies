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
	case DirUp:
		return elevio.MotorUp
	case DirDown:
		return elevio.MotorDown
	case DirStop:
		return elevio.MotorStop
	default:
		return elevio.MotorStop
	}
}

type Dirn int

const (
	DirUp   Dirn = 1
	DirDown Dirn = -1
	DirStop Dirn = 0
)

type ElevatorBehaviour int

const (
	ElevIdle     ElevatorBehaviour = 0
	ElevDoorOpen ElevatorBehaviour = 1
	ElevMoving   ElevatorBehaviour = 2
)

type Elevator struct {
	ID                 int
	Floor              int
	Dirn               Dirn
	Behaviour          ElevatorBehaviour
	Requests           [NumFloors][NumButtons]bool
	DoorOpenDuration_s float64
	Obstruction        bool
	MotorError         bool
}

// Elevator initializer function
func NewElevator() Elevator {
	var elevatorSetup Elevator = Elevator{
		ID:                 0,
		Floor:              -1,                            // Uninitialized floor
		Dirn:               DirStop,                       // Not moving
		Behaviour:          ElevIdle,                      // Idle state
		Requests:           [NumFloors][NumButtons]bool{}, // No requests initially
		DoorOpenDuration_s: 3.0,
		Obstruction:        false,
		MotorError:         false,
	}

	return elevatorSetup
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
	setElevatorObtruction(status)
	if status {
		StartObstrTimer()
	} else {
		StopObstrTimer()
	}
	fmt.Println("Obstruction status set to:", status)
}

func SetElevatorToValidStartPosition() {
	fmt.Println("Elevator initialized")
	for {
		if elevio.GetFloor() == -1 {
			elevio.SetMotorDirection(elevio.MotorDown)
		} else {
			elevio.SetMotorDirection(elevio.MotorStop)
			break
		}
		time.Sleep(_timerPollRate)

	}
	setAllLights(elevator)
}

// Fucntion that sets all requests and lights to false
func ClearAllRequests() {
	elevator.Requests = [NumFloors][NumButtons]bool{}
	setAllLights(elevator)
}
