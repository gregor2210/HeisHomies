package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
)

func SetAllLights() {
	// henter alle online heiser.
	// Lager en ny request matrise med alle knappetrykkene, så setter vi statusen på alle de
	var requests [fsm.NumFloors][fsm.NumButtons - 1]bool
	online_ids := Get_all_online_ids()

	for _, id := range online_ids {
		var req [fsm.NumFloors][fsm.NumButtons]bool // deafult false

		if id == ID {
			req = fsm.GetElevatorStruct().Requests
		} else {
			req = Get_worldview(id).Elevator.Requests
		}

		// checking hall up and down and copy if true
		for floor := 0; floor < fsm.NumFloors; floor++ {
			// For nr of floors

			if req[floor][0] {
				// Hall up == ture
				fmt.Println("UP")
				requests[floor][0] = true
			}

			if req[floor][1] {
				// Hall down == ture
				fmt.Println("DOWN")
				requests[floor][1] = true
			}
		}
	}

	fmt.Println("Requests:")
	for floor := 0; floor < fsm.NumFloors; floor++ {
		for btn := 0; btn < fsm.NumButtons-1; btn++ {
			if requests[floor][btn] {
				fmt.Printf("Floor %d, Button %d\n", floor, btn)
			}
		}
	}

	for floor := 0; floor < fsm.NumFloors; floor++ {
		for btn := 0; btn < fsm.NumButtons-1; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, requests[floor][btn])
		}
	}
}
