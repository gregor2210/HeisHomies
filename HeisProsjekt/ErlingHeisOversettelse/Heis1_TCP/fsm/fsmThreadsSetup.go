package fsm

import "fmt"

// Channel to receive timer TimeOut events
func FsmThreadsSetup() (chan bool, chan bool) {
	fmt.Println("Setting up FSM threads")
	timerTimeOutChan := make(chan bool)
	motorErrorChan = make(chan bool)

	go PollTimerTimeOut(timerTimeOutChan)
	go PollMotorTimerTimeOut(motorErrorChan)
	fmt.Println("FSM threads set up")

	return timerTimeOutChan, motorErrorChan
}
