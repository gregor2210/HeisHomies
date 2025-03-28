package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
)

// To store the last request and reduce network traffic
var lastRequest [fsm.NumFloors][fsm.NumButtons - 1]bool

func SetAllLights() {
	// Creates a new request matrix with all button presses, then updates their statuses
	var requests [fsm.NumFloors][fsm.NumButtons - 1]bool
	onlineIDs := GetAllOnlineIds()
	onlineIDs = append(onlineIDs, ID) // Include self
	for _, id := range onlineIDs {
		var req [fsm.NumFloors][fsm.NumButtons]bool // Default false

		if id == ID {
			req = fsm.GetElevatorStruct().Requests
		} else {
			req = GetWorldView(id).Elevator.Requests
		}

		// Check and copy hall buttons
		for floor := 0; floor < fsm.NumFloors; floor++ {

			// Hall up == true
			if req[floor][0] {
				requests[floor][0] = true
			}

			// Hall down == true
			if req[floor][1] {
				requests[floor][1] = true
			}
		}
	}

	// Skip if requests unchanged to save traffic
	if requests != lastRequest {

		for floor := 0; floor < fsm.NumFloors; floor++ {
			for btn := 0; btn < fsm.NumButtons-1; btn++ {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, requests[floor][btn])
			}
		}
	}

	lastRequest = requests
}
