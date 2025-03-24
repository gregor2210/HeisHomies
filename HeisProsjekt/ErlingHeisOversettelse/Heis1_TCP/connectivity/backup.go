package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
)

// Start backupprosess for dead elevator
func StartBackupProcess(deadElevID int) {

	deadWorldView := GetWorldView(deadElevID)
	// Extract requests from dead elevator
	deadRequests := deadWorldView.Elevator.Requests

	for i, floor := range deadRequests {
		if floor[0] {
			var button elevio.ButtonType = elevio.BtnHallUp
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request, true)

		}
		if floor[1] {
			var button elevio.ButtonType = elevio.BtnHallDown
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request, true)
		}

	}

}

// Start backupprosess for dead elevator
func StartMotorErrorBackupProcess() {
	fmt.Println("Starting backup prosess")

	deadRequests := fsm.GetElevatorStruct().Requests

	for i, floor := range deadRequests {
		if floor[0] {
			var button elevio.ButtonType = elevio.BtnHallUp
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request, false)

		}
		if floor[1] {
			var button elevio.ButtonType = elevio.BtnHallDown
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request, false)
		}

	}

}
