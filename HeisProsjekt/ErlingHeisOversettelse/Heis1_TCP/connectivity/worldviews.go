package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"sync"
)

var (
	world_view_backup       [NR_OF_ELEVATORS]Worldview_package
	world_view_backup_mutex sync.Mutex
)

func Store_worldview(id int, worldview Worldview_package) {
	world_view_backup_mutex.Lock()
	defer world_view_backup_mutex.Unlock()
	world_view_backup[id] = worldview

}

func Get_worldview(id int) Worldview_package {
	world_view_backup_mutex.Lock()
	defer world_view_backup_mutex.Unlock()
	return world_view_backup[id]
}

func Dose_order_exist(button_event elevio.ButtonEvent) bool {
	fmt.Println("Starting Dose order exist")
	floor := button_event.Floor
	var button int = int(button_event.Button) // 0 hallup, 1 halldown
	fmt.Println("Dose order exist button type: ", button, " Floor: ", floor)
	id_of_online_elevators := Get_all_online_ids()
	for _, id := range id_of_online_elevators {
		if id == ID {
			requests := fsm.GetElevatorStruct().Requests
			if requests[floor][button] {
				return true
			}

		} else {
			world_view := Get_worldview(id)
			requests := world_view.Elevator.Requests

			if requests[floor][button] {
				return true
			}
		}
	}

	return false
}
