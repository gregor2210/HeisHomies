package fsm

//Direction type. up = 1, down = 0
type Direction int

const (
	Up    Direction = 1
	Down  Direction = -1
	Still Direction = 0
)

//Elevator struct containing floor, moving direction and requests
//is used to keep track of the elevators state
//is basicly a elevator object
type Elevator struct {
	floor              int
	movingDirection    Direction
	moving_up_stopps   []int
	moving_dwon_stopps []int
}
