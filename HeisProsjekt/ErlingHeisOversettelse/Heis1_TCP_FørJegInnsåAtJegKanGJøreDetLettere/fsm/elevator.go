package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

// NumFloors and NumButtons are global variables
const NumFloors int = 4
const NumButtons int = 3

// Converting Dirn to MotorDirection
// For å få hvilken retning motoren fysisk skal gå basert på planlagt retning
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

// Direction type. up = 1, down = 0
type Dirn int

const (
	DirUp   Dirn = 1
	DirDown Dirn = -1
	DirStop Dirn = 0
)

// ElevatorBehaviour type. Idle = 0, DoorOpen = 1, Moving = 2
type ElevatorBehaviour int

const (
	ElevIdle     ElevatorBehaviour = 0
	ElevDoorOpen ElevatorBehaviour = 1
	ElevMoving   ElevatorBehaviour = 2
)

// Elevator struct containing floor, moving direction and requests
// is used to keep track of the elevators state
// is basicly a elevator object
type Elevator struct {
	ID        int
	Floor     int
	Dirn      Dirn
	Behaviour ElevatorBehaviour
	//Buttons in hall and cab x=floor y=button
	Requests           [NumFloors][NumButtons]bool
	DoorOpenDuration_s float64
	Obstruction        bool
}

// Elevator initializer function
func NewElevator() Elevator {
	var elevatorSetup Elevator = Elevator{
		ID:                 0,
		Floor:              -1,           // Uninitialized floor
		Dirn:               DirStop,      // Not moving
		Behaviour:          ElevIdle,     // Idle state
		Requests:           [4][3]bool{}, // No requests initially
		DoorOpenDuration_s: 3.0,
		Obstruction:        false, // Default door open duration
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

// Function to set the obstruction status of the elevator
func SetObsructionStatus(status bool) {
	elevator.Obstruction = status
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
		time.Sleep(_pollRate)

	}
}
