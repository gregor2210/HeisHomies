package fsm

// Direction type. up = 1, down = 0
type Direction int

const (
	Up    Direction = 1
	Down  Direction = -1
	Still Direction = 0
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
	floor           int
	movingDirection Direction
	behaviour       ElevatorBehaviour
	//Buttons in hall and cab x=floor y=button
	requests           [4][3]bool
	doorOpenDuration_s float64
}
