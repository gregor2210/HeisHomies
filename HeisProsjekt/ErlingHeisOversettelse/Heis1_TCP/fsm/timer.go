package fsm

import (
	"fmt"
	"time"
)

// IF THIS CHANGES, REMEMBER TO UPDATE IT IN ELEVATOR:IO.GO AS WELL
const _pollRate = 20 * time.Millisecond

var timerEndTime time.Time
var timerActive bool

// Polls the timer and signals on TimeOut
func PollTimerTimeOut(receiver chan<- bool) {
	for {
		time.Sleep(_pollRate) // Poll rate, adjust as needed
		if TimerTimedOut() {
			receiver <- true
		}
	}
}

// Current time
func getWallTime() time.Time {
	return time.Now()
}

// Starts a timer with a given duration
func TimerStart(duration float64) {
	fmt.Println("Timer started, for:", duration)
	timerEndTime = getWallTime().Add(time.Duration(duration * float64(time.Second)))
	timerActive = true
}

// Stops running timer
func TimerStop() {
	timerActive = false
}

// Returns true if timer expired and no obstruction
func TimerTimedOut() bool {
	return timerActive && getWallTime().After(timerEndTime) && !(elevator.Obstruction)
}
