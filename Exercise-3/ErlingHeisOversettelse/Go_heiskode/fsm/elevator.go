package fsm

import (
	"Driver-go/elevio"
)

// NumFloors and NumButtons are global variables
var NumFloors int = 4
var NumButtons int = 3

// Converting Dirn to MotorDirection
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

// Direction type. up = 1, down = 0
type Dirn int

const (
	D_Up   Dirn = 1
	D_Down Dirn = -1
	D_Stop Dirn = 0
)

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
	dirn      Dirn
	behaviour ElevatorBehaviour
	//Buttons in hall and cab x=floor y=button
	requests           [4][3]bool
	doorOpenDuration_s float64
}

// Elevator initializer function
func NewElevator() Elevator {
	return Elevator{
		floor:              -1,           // Uninitialized floor
		dirn:               D_Stop,       // Not moving
		behaviour:          EB_Idle,      // Idle state
		requests:           [4][3]bool{}, // No requests initially
		doorOpenDuration_s: 3.0,          // Default door open duration
	}
}
