package fsm

import (
	"fmt"
	"time"
)

var (
	timerTimeOutChan = make(chan bool)
	timerEndTime     time.Time
	timerActive      bool
)

// Polls the timer and signals on TimeOut
func PollTimerTimeOut() {
	for {
		time.Sleep(_timerPollRate) 
		if TimerTimedOut() {
			timerTimeOutChan <- true
		}
	}
}

// Current time
func getWallTime() time.Time {
	return time.Now()
}

// Starts a timer with a given duration
func TimerDoorStart(duration float64) {
	fmt.Println("Timer started, for:", duration)
	timerEndTime = getWallTime().Add(time.Duration(duration * float64(time.Second)))
	timerActive = true
}

// Stops running timer
func TimerStop() {
	timerActive = false
}

// Returns true if timer expires and no obstruction
func TimerTimedOut() bool {
	return timerActive && getWallTime().After(timerEndTime) && !(elevator.Obstruction)
}
