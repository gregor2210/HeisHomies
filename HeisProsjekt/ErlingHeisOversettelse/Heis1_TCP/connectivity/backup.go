package connectivity

import "Driver-go/elevio"

func Start_backup_process(dead_elevator_id int) {
	dead_worldview := Get_worldview(dead_elevator_id)

	//extract current requests.
	//var new_requests []elevio.ButtonEvent
	dead_requests := dead_worldview.Elevator.Requests
	for i, floor := range dead_requests {
		if floor[0] {
			var button elevio.ButtonType = elevio.BT_HallUp
			request := elevio.ButtonEvent{Floor: i, Button: button}
			//new_requests = append(new_requests, elevio.ButtonEvent{Floor: i, Button: button})
			New_order(request)

		}
		if floor[1] {
			var button elevio.ButtonType = elevio.BT_HallDown
			request := elevio.ButtonEvent{Floor: i, Button: button}
			New_order(request)

			//new_requests = append(new_requests, elevio.ButtonEvent{Floor: i, Button: button})
		}

	}

}
