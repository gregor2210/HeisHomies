package fsm

// Channel to receive timer TimeOut events
func FsmThreadsSetup() chan bool {
	timerTimeOutChan := make(chan bool)
	go PollTimerTimeOut(timerTimeOutChan)

	return timerTimeOutChan
}
