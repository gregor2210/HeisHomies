package fsm

import (
	"Driver-go/elevio"
)

// NumFloors and NumButtons are global variables
var NumFloors int = 4
var NumButtons int = 3

// ElevatorBehaviour type. Idle = 0, DoorOpen = 1, Moving = 2
type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2
)

// Elevator struct containing floor, moving direction and requests
// is used to keep track of the elevators state
// is basicly a elevator object
type Elevator struct {
	floor     int
	dirn      elevio.MotorDirection
	behaviour ElevatorBehaviour
	//Buttons in hall and cab x=floor y=button
	requests           [4][3]bool
	doorOpenDuration_s float64
}

// Elevator initializer function
func NewElevator() Elevator {
	return Elevator{
		floor:              -1,             // Uninitialized floor
		dirn:               elevio.MD_Stop, // Not moving
		behaviour:          EB_Idle,        // Idle state
		requests:           [4][3]bool{},   // No requests initially
		doorOpenDuration_s: 3.0,            // Default door open duration
	}
}
