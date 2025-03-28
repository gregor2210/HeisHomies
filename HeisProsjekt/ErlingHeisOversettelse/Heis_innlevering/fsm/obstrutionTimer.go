package fsm

import (
	"fmt"
	"time"
)

var (
	obstrErrorChan    = make(chan bool)
	obstrTimerEndTime time.Time
	obstrTimerActive  bool
)

// Polls the timer and signals on TimeOut
func PollObstrTimerTimeOut() {
	for {
		time.Sleep(_timerPollRate) 
		if ObstrTimerTimedOut() {
			fmt.Println("Obs error timer timed out")
			obstrErrorChan <- true
			obstrTimerActive = false
		}
	}
}

// Starts obstruction timer
func StartObstrTimer() {
	var duration float64 = _obstrErrorDuration
	fmt.Println("Obs timer started, for:", duration)
	obstrTimerEndTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
	obstrTimerActive = true

}

// Stops running timer
func StopObstrTimer() {
	fmt.Println("Obs error timer stoped")
	obstrTimerActive = false
}

// Returns true if timer expires and no obstruction
func ObstrTimerTimedOut() bool {
	return obstrTimerActive && time.Now().After(obstrTimerEndTime)
}
