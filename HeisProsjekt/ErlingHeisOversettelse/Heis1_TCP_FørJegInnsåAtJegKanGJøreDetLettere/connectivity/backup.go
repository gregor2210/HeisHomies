package connectivity

import "Driver-go/elevio"

func Start_backup_prosess(deadElevID int) {
	deadWorldView := GetWorldView(deadElevID)

	//extract current requests.
	var new_requests []elevio.ButtonEvent
	deadRequests := deadWorldView.Elevator.Requests
	for i, floor := range deadRequests {
		if floor[0] {
			var button elevio.ButtonType = elevio.BtnHallUp
			new_requests = append(new_requests, elevio.ButtonEvent{Floor: i, Button: button})

		} else if floor[1] {
			var button elevio.ButtonType = elevio.BtnHallDown
			new_requests = append(new_requests, elevio.ButtonEvent{Floor: i, Button: button})
		}

	}

}
