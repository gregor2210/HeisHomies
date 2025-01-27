package fsm

import (
	"time"
)

var timerEndTime time.Time
var timerActive bool

func getWallTime() time.Time {
	return time.Now()
}

func TimerStart(duration time.Duration) {
	timerEndTime = getWallTime().Add(duration)
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimerTimedOut() bool {
	return timerActive && getWallTime().After(timerEndTime)
}