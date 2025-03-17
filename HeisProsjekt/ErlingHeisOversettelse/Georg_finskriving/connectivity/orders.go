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

func Calculate_priority_value(button_event elevio.ButtonEvent, elevator fsm.Elevator) int {
	request_floor := button_event.Floor
	//request_button := button_event.Button

	//button point dir. minus is down
	//requst_button_point_dir := -1
	//if int(request_button) == 0 {
	//requst_button_point_dir = 1
	//}
	//elevator := GetElevatorStruct()
	NumFloors_minus_1 := fsm.NumFloors - 1
	//Calculate how much this elevator wants this request.
	priority_value := 2 * 10 * NumFloors_minus_1 // max value

	//DÅRLIG VERSJON, TAR ABSOLUTT AVSTAND
	delta_floor := request_floor - elevator.Floor
	//fmt.Println("---------------------------------------------------")
	//fmt.Println("request floor: ", request_floor)
	//fmt.Println("elevator floor: ", elevator.Floor)
	//fmt.Println("Delta floor: ", delta_floor)
	//fmt.Println("Priority value :", priority_value)

	priority_value = priority_value - int(math.Abs(float64(delta_floor)))*10

	//FUGERTE IKKE HELT
	/*
		//antall etasjer unna gir -10 poeng
		//delta_floor gir minus value if requested floor is below
		delta_floor := request_floor - elevator.Floor

		//if elevator dosen ot have a moving dirn
		if int(elevator.Dirn) == 0 {
			priority_value = priority_value - int(math.Abs(float64(delta_floor)))*10

			//HVis heis beveger seg mot request, og request peker i samme retning som heisens bevegelse
		} else if math.Copysign(1, float64(delta_floor)) == math.Copysign(1, float64(elevator.Dirn)) {
			// delta_floor have the same sign as elv direction. Meaning it is going towards the request)

			if math.Copysign(1, float64(requst_button_point_dir)) == math.Copysign(1, float64(elevator.Dirn)) {
				//Elevator moves towards request and in same direction as request
				priority_value = priority_value - int(math.Abs(float64(priority_value)-float64(delta_floor)))*10

			} else {
				//Elevator moves toward request, but request is in oposit riection
				if int(elevator.Dirn) < 0 {
					//elevator moves downward
					wortcase_down := elevator.Floor + request_floor //goes to bottom, and turns direciton and moves up
					priority_value = priority_value - wortcase_down*10
				} else {
					//elevator moves up
					worstcase_up := (NumFloors_minus_1 - elevator.Floor) + (NumFloors_minus_1 - request_floor)
					priority_value = priority_value - worstcase_up*10
				}
			}

		} else {
			// Elevator is not going towards the request

			if math.Copysign(1, float64(requst_button_point_dir)) == math.Copysign(1, float64(elevator.Dirn)) {
				//Elevator will not point in same drection after turn at a worstcase top

			} else {
				//Elevator moves toward request, but request is in oposit riection
				if int(elevator.Dirn) < 0 {
					//elevator moves downward
					nr_of_floors_traveled := elevator.Floor + NumFloors_minus_1 + request_floor
					priority_value = priority_value - nr_of_floors_traveled*10
				} else {
					//elevator moves up
					nr_of_floors_traveled := (NumFloors_minus_1 - elevator.Floor) + NumFloors_minus_1 + request_floor
					priority_value = priority_value - nr_of_floors_traveled*10

				}
			}

		}*/
	fmt.Println("Priority value:", priority_value)
	//fmt.Println("---------------------------------------------------")
	return priority_value
}

func New_order(button_event elevio.ButtonEvent) {
	fmt.Println("New Order")

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

	//Sorting priorityvalue_to_sort in decending order
	sort.Sort(sort.Reverse(sort.IntSlice(priorityvalue_to_sort)))

	//Finding the elv id with highest priority value. and trying to send order to that elevator
	for _, priority_value := range priorityvalue_to_sort {
		//find id or elevator that will get
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
