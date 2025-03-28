package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

var (
	motorErrorChan    = make(chan bool)
	motorTimerEndTime time.Time
	motorTimerActive  bool
)

// Polls the timer and signals on TimeOut
func PollMotorTimerTimeOut() {
	for {
		time.Sleep(_timerPollRate) 
		if MotorTimerTimedOut() {
			fmt.Println("Motor error timer timed out")
			motorErrorChan <- true
			motorTimerActive = false
		}
	}
}

// Starts motor error timer if behaviour is not MotorStop
func StartMotorErrorTimer(elv Elevator) {
	if elv.Behaviour != elevio.MotorStop {
		var duration float64 = _motorErrorDuration
		fmt.Println("Motor timer started, for:", duration)
		motorTimerEndTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
		motorTimerActive = true
	}

}

// Stops running timer
func StopMotorErrorTimer() {
	fmt.Println("Motor error timer stoped")
	motorTimerActive = false
}

// Returns true if timer expires and no obstruction
func MotorTimerTimedOut() bool {
	return motorTimerActive && time.Now().After(motorTimerEndTime)
}
