package fsm

import (
	"fmt"
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
