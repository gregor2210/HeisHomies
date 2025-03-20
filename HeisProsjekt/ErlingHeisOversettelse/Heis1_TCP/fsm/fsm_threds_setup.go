package fsm

func Fsm_threds_setup() chan bool {
	// Channel to receive timer timeout events
	timerTimeoutChan := make(chan bool)
	go PollTimerTimeout(timerTimeoutChan)

	return timerTimeoutChan
}
