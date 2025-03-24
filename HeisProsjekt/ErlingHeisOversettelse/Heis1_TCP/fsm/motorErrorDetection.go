package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

// IF THIS CHANGES, REMEMBER TO UPDATE IT IN ELEVATOR:IO.GO AS WELL
const _motorPullRate = 20 * time.Millisecond

var (
	motorErrorChan    chan bool
	motorTimerEndTime time.Time
	motorTimerActive  bool
)

// Polls the timer and signals on TimeOut
func MotorTimerTimeOut(receiver chan<- bool) {
	for {
		time.Sleep(_motorPullRate) // Poll rate, adjust as needed
		if MotorTimerTimedOut() {
			//
			fmt.Println("Motor error timer time out")
			motorErrorChan <- true
		}
	}
}

// Starts motor error timer if behaviour is not MotorStop
func StartMotorErrorTimer(elv Elevator) {
	if elv.Behaviour != elevio.MotorStop {
		var duration float64 = 10 // 10sec
		fmt.Println("Motor timer started, for:", duration)
		motorTimerEndTime = getWallTime().Add(time.Duration(duration * float64(time.Second)))
		motorTimerActive = true
	}

}

// Stops running timer
func StopMotorErrorTimer() {
	fmt.Println("Motor error timer stoped")
	motorTimerActive = false
}

// Returns true if timer expired and no obstruction
func MotorTimerTimedOut() bool {
	return motorTimerActive && getWallTime().After(motorTimerEndTime)
}
