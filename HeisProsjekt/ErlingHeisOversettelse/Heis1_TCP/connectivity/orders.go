package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"math"
	"sort"
)

type OrderRequests struct {
	Original_elevator               int
	Responder_elevator              int
	Unique_ID                       string
	Button_event                    elevio.ButtonEvent
	Elevator_priority_value_respond int
}

func PrintOrderRequest(orders []OrderRequests) {
	fmt.Println("Pending Orders:")
	for _, order := range orders {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.Original_elevator, order.Responder_elevator, order.Unique_ID, order.Elevator_priority_value_respond)
	}
}

func Calculate_priority_value(button_event elevio.ButtonEvent, e fsm.Elevator) int {
	// Calculating the priority value.
	// Larger value == higher order priority.
	request_floor := button_event.Floor

	NumFloors_minus_1 := fsm.NumFloors - 1
	//Calculate how much this elevator wants this request.
	priority_value := 2 * 10 * NumFloors_minus_1 // max value

	//DÅRLIG VERSJON, TAR ABSOLUTT AVSTAND
	delta_floor := request_floor - e.Floor

	sub_val := int(math.Abs(float64(delta_floor))) * 10
	//if elevator dosen ot have a moving dirn
	if int(e.Dirn) == 0 {
		// Elevator ha no order
		//-

	} else if delta_floor < 0 && int(e.Dirn) < 0 {
		// Elevator moves down toward the button, down

		if button_event.Button == elevio.BT_HallDown {
			// Button moves in the same direction as the elevator: down
			//-

		} else {
			//  Button points in the opposite direction of the elevator: up
			sub_val += button_event.Floor * 2 * 10

		}

	} else if delta_floor > 0 && int(e.Dirn) < 0 {
		// Elevator moves down, away from the button
		sub_val += 2 * e.Floor * 10

		if button_event.Button == elevio.BT_HallUp {
			// Button points up
			//-
		} else {
			// Button points down
			sub_val += (NumFloors_minus_1 - button_event.Floor) * 2 * 10

		}

	} else if delta_floor > 0 && int(e.Dirn) > 0 {
		// Elevator moves up toward the button, up

		if button_event.Button == elevio.BT_HallUp {
			// Button points in the same direction: up
			//-

		} else {
			// Knapp peker motsatt vei NED
			sub_val += (NumFloors_minus_1 - button_event.Floor) * 2 * 10
		}
	} else if delta_floor < 0 && int(e.Dirn) > 0 {
		// Elevator moves up, away from the button
		sub_val += (NumFloors_minus_1 - e.Floor) * 2 * 10

		if button_event.Button == elevio.BT_HallUp {
			// Button points up
			sub_val += button_event.Floor * 2 * 10
		} else {
			// Button points down
			//-
		}
	}

	priority_value -= sub_val

	//fmt.Println("Priority value:", priority_value)
	return priority_value
}

func New_order(button_event elevio.ButtonEvent) {
	// Figure out who should take which order.
	// Sends the order to the selected elevator

	if Dose_order_exist(button_event) {
		fmt.Println("Order allready exist")
		return
	}

	var priority_value_id_index [NR_OF_ELEVATORS]int
	var priorityvalue_to_sort []int
	online_elevator_id := Get_all_online_ids()
	for _, id := range online_elevator_id {
		//Gets list of all priority values
		var elevator fsm.Elevator
		if id == ID {
			elevator = fsm.GetElevatorStruct()

		} else {
			elevator = Get_worldview(id).Elevator

		}

		priority_value := Calculate_priority_value(button_event, elevator)
		priority_value_id_index[id] = priority_value
		priorityvalue_to_sort = append(priorityvalue_to_sort, priority_value)
	}

	// Sorting priorityvalue_to_sort in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(priorityvalue_to_sort)))

	// Finding the elevator ID with the highest priority value and attempting to send the order to that elevator
	for _, priority_value := range priorityvalue_to_sort {
		// Find the elevator ID that will receive the order
		id_of_elevator_that_will_get_order := ID
		for i, v := range priority_value_id_index {
			if v == priority_value {
				id_of_elevator_that_will_get_order = i
				break
			}
		}

		if id_of_elevator_that_will_get_order == ID {
			//Send reqeust to self
			fsm.Fsm_onRequestButtonPress(button_event.Floor, button_event.Button) // Ikke så fint at dnne er her
			break

		} else if Send_order_to_spesific_elevator(id_of_elevator_that_will_get_order, button_event) {
			//Try to send order to elevator with id
			fmt.Println("ID: ", id_of_elevator_that_will_get_order, "Got the order!")
			break
		}

	}

}
