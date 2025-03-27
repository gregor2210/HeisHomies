package fsm

import "time"

const (
	_timerPollRate              = 20 * time.Millisecond
	_motorErrorDuration float64 = 20 // sec
	_obstrErrorDuration float64 = 20 // sec
	NumFloors           int     = 4
	NumButtons          int     = 3
)
