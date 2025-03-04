package connectivity

import (
	"Driver-go/elevio"
	"fmt"
	"sort"
	"sync"
	"time"
)

type OrderRequests struct {
	Original_elevator        int
	Responder_elevator       int
	Time                     string
	Button_event             elevio.ButtonEvent
	Elevator_priority_values [NR_OF_ELEVATORS]int
}

type SingleResponse struct {
	receved_from_id int
	priority_value  int
}

var (
	pending_orders []OrderRequests
	order_response []OrderRequests

	my_request_responses_chans = make(map[string]chan SingleResponse)
	pending_orders_mutex       sync.Mutex
)

func New_order(floor int, button int) {
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()

	butten_type := elevio.ButtonType(button)
	button_event := elevio.ButtonEvent{Floor: floor, Button: butten_type}

	now := time.Now()
	now_str := now.Format("2006/01/02 15:04:05.999999999")

	new_order := OrderRequests{Original_elevator: ID, Time: now_str, Button_event: button_event}

	pending_orders = append(pending_orders, new_order)

	//setting up listening to repsond go rutine

	fmt.Println("now_str", now_str)
	response_chan := make(chan SingleResponse)
	//add it to dict to later be able to send SingleResponse to the corrrect go rutine based on the now_str
	my_request_responses_chans[now_str] = response_chan

	priority_value := calculate_priority_value()
	go Wait_for_order_responses(now_str, response_chan, priority_value, Get_all_online_ids())

}

func Wait_for_order_responses(now_string string, singel_response_chan chan SingleResponse, this_elevators_priority_value int, online_elevators []int, order elevio.ButtonEvent) {
	var responses [NR_OF_ELEVATORS]int
	//responses_to_sort is a variabel i can sort later the nuse it to fin the indexes.
	//it makes it possible to use builtinn sort func
	var responses_to_sort []int
	response_counter := 0

	responses[ID] = this_elevators_priority_value
	responses_to_sort = append(responses_to_sort, this_elevators_priority_value)
	response_counter++

	//Timeoout will kick in after given miliseconds. The elevator will give the respons to the highest values.
	var timeout <-chan time.Time
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	timeout = ticker.C

loop:
	for {
		select {
		case single_response := <-singel_response_chan:
			if responses[single_response.receved_from_id] == 0 {
				responses[single_response.receved_from_id] = single_response.priority_value
				responses_to_sort = append(responses_to_sort, single_response.priority_value)
				response_counter++
			}
			if response_counter == len(online_elevators) {
				delete(my_request_responses_chans, now_string)
				break loop
			}

		case <-timeout:
			delete(my_request_responses_chans, now_string)
			break loop
		}

	}

	//Choose who will get the order. find highest value and then the correspoinding index
	sort.Sort(sort.Reverse(sort.IntSlice(responses_to_sort))) // Sort descending
	for _, priority_value := range responses_to_sort {

		//find id or elevator that will get
		id_of_elevator_that_will_get_order := ID
		for i, v := range responses {
			if v == priority_value {
				id_of_elevator_that_will_get_order = i
				break
			}
		}

		if id_of_elevator_that_will_get_order == ID {
			//Send reqeust to self
			break

		} else if Send_order_to_spesific_elevator(id_of_elevator_that_will_get_order, order) {
			break
		}

	}
	//remove the order from the pending order!

}

func Get_order_requests() []OrderRequests {
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()
	return pending_orders
}

func Get_order_respons() []OrderRequests {
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()
	return order_response
}

func calculate_priority_value() int {
	return 4
}

func Receved_order_requests(received_order_requests []OrderRequests) {
	//Receves the incomming order_request_array
	//Removes any respons to a request that is not in request with the same id/signature.
	// Itterate through it and chech if we have responded to it before. by comparing the 3 first parameters
	//if not, calc priority_value and add an order response.
	//if yes, skip

	// For elm in received_order_requests with id of the sender of these requests, check if there is an

	//Deleting responses that are not received_order_requests with the same sender id
	//Iterate over each order in order_response
	for j := 0; j < len(order_response); j++ {
		order := order_response[j]
		// Check if the order is in received_order_requests
		found := false
		for _, received := range received_order_requests {
			if order.Original_elevator == received.Original_elevator && order.Time == received.Time {
				found = true
				break
			}
		}
		// If the order is not in received_order_requests, remove it from order_response
		if !found {
			//Dette er sånn man fjærner element j fra en liste
			order_response = append(order_response[:j], order_response[j+1:]...)
			j-- // Decrement j to account for the removed element
		}
	}

	//Generating new response
	for _, receved_order := range received_order_requests {
		//Deleating responses that are no receved_order_requests with same sender id
		//for each receved_order
		for _, order := range order_response {
			//is that order in order_response
			if order.Original_elevator == receved_order.Original_elevator && order.Time == receved_order.Time {
				// The order has already been responded to, so skip it
				continue
			}
			// Calculate priority_value and add an order response
			priority_value := calculate_priority_value()
			receved_order.Responder_elevator = ID
			receved_order.Elevator_priority_values[ID] = priority_value
			order_response = append(order_response, receved_order)
			// ...
		}
	}

}

func Receved_order_response(received_order_response []OrderRequests) {
	for _, order := range received_order_response {
		if order.Original_elevator == ID {
			// Do something here
		}
	}
}
