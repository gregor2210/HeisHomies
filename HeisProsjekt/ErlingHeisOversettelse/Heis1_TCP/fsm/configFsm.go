package fsm

import "time"

// IF THIS CHANGES, REMEMBER TO UPDATE IT IN ELEVATOR:IO.GO AS WELL
const (
	_timerPollRate              = 20 * time.Millisecond
	_motorErrorDuration float64 = 20 // sec
	_obstrErrorDuration float64 = 20 // sec
	NumFloors           int     = 4
	NumButtons          int     = 3
)
