package connectivity

import (
	"Driver-go/elevio"
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
			NewOrder(request)

		}
		if floor[1] {
			var button elevio.ButtonType = elevio.BtnHallDown
			request := elevio.ButtonEvent{Floor: i, Button: button}
			NewOrder(request)
		}

	}

}

// Start backupprosess for broken elevator
func DisableComunicaton() {
	SetSelfOffline()
	CloseAllConnections()
}
