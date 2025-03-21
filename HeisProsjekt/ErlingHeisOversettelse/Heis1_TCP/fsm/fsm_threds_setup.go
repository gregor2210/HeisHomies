package fsm

// Channel to receive timer timeout events
func Fsm_threds_setup() chan bool {
	timerTimeoutChan := make(chan bool)
	go PollTimerTimeout(timerTimeoutChan)

	return timerTimeoutChan
}
