package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
)

// To store the last request
// For reducing network traffic
var lastRequest [fsm.NumFloors][fsm.NumButtons - 1]bool

// Set the hall lights on button panel, using worldview_backups and self
func SetAllLights() {
	// Retrieves all online elevators.
	// Creates a new request matrix with all button presses, then updates their statuses
	var requests [fsm.NumFloors][fsm.NumButtons - 1]bool
	onlineIDs := GetAllOnlineIds()
	onlineIDs = append(onlineIDs, ID) // Should allways include self
	for _, id := range onlineIDs {
		var req [fsm.NumFloors][fsm.NumButtons]bool // deafult false

		if id == ID {
			req = fsm.GetElevatorStruct().Requests
		} else {
			req = GetWorldView(id).Elevator.Requests
		}

		// Checking hall up and down buttons and copying if true
		for floor := 0; floor < fsm.NumFloors; floor++ {
			// For nr of floors

			if req[floor][0] {
				// Hall up == ture
				//fmt.Println("UP")
				requests[floor][0] = true
			}

			if req[floor][1] {
				// Hall down == ture
				//fmt.Println("DOWN")
				requests[floor][1] = true
			}
		}
	}

	// Printing for debugging
	/*fmt.Println("Requests:")
	for floor := 0; floor < fsm.NumFloors; floor++ {
		for btn := 0; btn < fsm.NumButtons-1; btn++ {
			if requests[floor][btn] {
				fmt.Printf("Floor %d, Button %d\n", floor, btn)
			}
		}
	}
	*/
	if requests != lastRequest {
		// If the requests are the same as before, there is no need to use network traffic to set the same states
		for floor := 0; floor < fsm.NumFloors; floor++ {
			for btn := 0; btn < fsm.NumButtons-1; btn++ {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, requests[floor][btn])
			}
		}
	}

	lastRequest = requests
}
