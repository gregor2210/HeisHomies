package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
)

// Start backupprosess for dead elevator
func StartBackupProcess(deadElevID int) {
	fmt.Println("Starting backup prosess")
	var deadRequests [fsm.NumFloors][fsm.NumButtons]bool
	if deadElevID == ID {
		deadRequests = fsm.GetElevatorStruct().Requests
	} else {
		deadWorldView := GetWorldView(deadElevID)
		// Extract requests from dead elevator
		deadRequests = deadWorldView.Elevator.Requests

	}
	for i, floor := range deadRequests {
		if floor[0] {
			var button elevio.ButtonType = elevio.BtnHallUp
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request)

		}
		if floor[1] {
			var button elevio.ButtonType = elevio.BtnHallDown
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request)
		}

	}

}
