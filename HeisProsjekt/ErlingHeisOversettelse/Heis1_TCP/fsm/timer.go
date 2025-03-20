package fsm

import (
	"fmt"
	"time"
)

// IF THIS CHANGES, REMEMBER TO UPDATE IT IN ELEVATOR:IO.GO AS WELL
const _pollRate = 20 * time.Millisecond

var timerEndTime time.Time
var timerActive bool

// Function that polls the timer at a given rate. Sends a signal when the timer has timed out
func PollTimerTimeout(receiver chan<- bool) {
	for {
		time.Sleep(_pollRate) // Poll rate, adjust as needed
		if TimerTimedOut() {
			receiver <- true // Send a signal when the timer has timed out
		}
	}
}

// Returns time Now
func getWallTime() time.Time {
	return time.Now()
}

// starts a timer with a given duration
func TimerStart(duration float64) {
	fmt.Println("Timer started, for:", duration)
	timerEndTime = getWallTime().Add(time.Duration(duration * float64(time.Second)))
	timerActive = true
}

// Stops running timer
func TimerStop() {
	timerActive = false
}

// Returns bool if timer is timedout or not.
func TimerTimedOut() bool {
	return timerActive && getWallTime().After(timerEndTime) && !(elevator.Obstruction)
}
