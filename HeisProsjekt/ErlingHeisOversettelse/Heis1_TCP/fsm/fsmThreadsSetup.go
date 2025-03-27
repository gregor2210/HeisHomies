package fsm

import (
	"fmt"
	"time"
)

// IF THIS CHANGES, REMEMBER TO UPDATE IT IN ELEVATOR:IO.GO AS WELL
const (
	_timerPollRate              = 20 * time.Millisecond
	_motorErrorDuration float64 = 20 // sec
	_obstrErrorDuration float64 = 20 // sec
	NumFloors           int     = 7
	NumButtons          int     = 3
)

// Channel to receive timer timeout events
func FsmThreadsSetup() (chan bool, chan bool, chan bool) {
	fmt.Println("Setting up FSM threads")

	// Load cab requests from file
	loadCabRequests(&elevator) // elevator is in fsm.go

	go PollTimerTimeOut()
	go PollMotorTimerTimeOut()
	go PollObstrTimerTimeOut()
	fmt.Println("FSM threads set up")

	return timerTimeOutChan, motorErrorChan, obstrErrorChan
}
