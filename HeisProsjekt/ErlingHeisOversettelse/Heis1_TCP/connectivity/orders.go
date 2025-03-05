package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"sort"
	"sync"
	"time"
)

type OrderRequests struct {
	Original_elevator               int
	Responder_elevator              int
	Time                            string
	Button_event                    elevio.ButtonEvent
	Elevator_priority_value_respond int
}

type SingleResponse struct {
	receved_from_id int
	priority_value  int
}

type DoneProcessedOrder struct {
	Responses        [NR_OF_ELEVATORS]int
	Sorted_responses []int
	Order            elevio.ButtonEvent
}

var (
	pending_orders []OrderRequests
	order_response []OrderRequests

	my_request_responses_chans = make(map[string]chan SingleResponse)

	pending_orders_mutex sync.Mutex
	order_response_mutex sync.Mutex

	order_to_send_chan chan DoneProcessedOrder
)

func Order_setup(order_to_send_chan_ chan DoneProcessedOrder) {
	order_to_send_chan = order_to_send_chan_
}

func Get_pending_orders() []OrderRequests {
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()
	return pending_orders
}

func Get_order_respons() []OrderRequests {
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()
	return order_response
}

func New_order(button_event elevio.ButtonEvent, priority_value int) {
	fmt.Println("New Order")
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()

	//Setting up the new order
	now := time.Now()
	now_str := now.Format("2006/01/02 15:04:05.999999999")
	fmt.Println("now_str", now_str)
	new_order := OrderRequests{Original_elevator: ID, Time: now_str, Button_event: button_event}
	pending_orders = append(pending_orders, new_order)

	//setting up listening to repsond go rutine
	response_chan := make(chan SingleResponse)
	//add it to dict to later be able to send SingleResponse to the corrrect go rutine based on the now_str
	my_request_responses_chans[now_str] = response_chan

	//priority_value := calculate_priority_value()
	go Wait_for_order_responses(now_str, response_chan, priority_value, Get_all_online_ids(), button_event)
	fmt.Println("New Order added and waiting response")
}

func Wait_for_order_responses(now_string string, singel_response_chan chan SingleResponse, this_elevators_priority_value int, online_elevators []int, order elevio.ButtonEvent) {
	fmt.Println("Started go rutine Wait_for_order_responses")
	var responses [NR_OF_ELEVATORS]int
	//responses_to_sort is a variabel that can be sort later, the use it to find the indexe/id of elevator to get the order
	//it makes it possible to use builtinn sort func
	var responses_to_sort []int
	response_counter := 0

	responses[ID] = this_elevators_priority_value
	responses_to_sort = append(responses_to_sort, this_elevators_priority_value)
	response_counter++

	//Timeoout will kick in after given miliseconds. The elevator will give the respons to the highest values.
	var timeout <-chan time.Time
	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	timeout = ticker.C

loop:
	for {
		select {
		case single_response := <-singel_response_chan:
			if responses[single_response.receved_from_id] == 0 {
				fmt.Println("Got a new response, added, from", single_response.receved_from_id)
				responses[single_response.receved_from_id] = single_response.priority_value
				responses_to_sort = append(responses_to_sort, single_response.priority_value)
				response_counter++
			}
			if response_counter == len(online_elevators) {
				fmt.Println("All responses receved")
				delete(my_request_responses_chans, now_string)
				break loop
			}

		case <-timeout:
			delete(my_request_responses_chans, now_string)
			break loop
		}

	}
	fmt.Println("wait for respond loop finished!")
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()
	//remove the order from the pending order!
	for i, v := range pending_orders {
		if v.Time == now_string {
			//removing index i from pending order
			pending_orders = append(pending_orders[:i], pending_orders[i+1:]...)
		}
	}

	//Choose who will get the order. find highest value and then the correspoinding index
	sort.Sort(sort.Reverse(sort.IntSlice(responses_to_sort))) // Sort descending
	done_processed_order := DoneProcessedOrder{responses, responses_to_sort, order}
	order_to_send_chan <- done_processed_order
}

func SendOrderToSpesificElevator(done_processed_order DoneProcessedOrder) bool {
	order := done_processed_order.Order
	responses := done_processed_order.Responses
	sorted_responses := done_processed_order.Sorted_responses

	return_bool := false

	for _, priority_value := range sorted_responses {

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
			return_bool = false
			break

		} else if Send_order_to_spesific_elevator(id_of_elevator_that_will_get_order, order) {
			fmt.Println("ID: ", id_of_elevator_that_will_get_order, "Got the order!")
			return_bool = true
			break
		}

	}
	return return_bool

}

func Receved_order_requests(received_order_requests []OrderRequests, received_from_id int) {
	fmt.Println("Received order request----------------------------------")
	fmt.Println("Pending Orders:")
	for _, order := range pending_orders {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.Original_elevator, order.Responder_elevator, order.Time, order.Elevator_priority_value_respond)
	}

	fmt.Println("\nOrder Response:")
	for _, order := range order_response {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.Original_elevator, order.Responder_elevator, order.Time, order.Elevator_priority_value_respond)
	}
	fmt.Println("\nRECEVED Order Response:")
	for _, order := range received_order_requests {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.Original_elevator, order.Responder_elevator, order.Time, order.Elevator_priority_value_respond)
	}
	order_response_mutex.Lock()
	defer order_response_mutex.Unlock()
	//Receves the incomming order_request_array
	//Removes any respons to a request that is not in request with the same id/signature.
	// Itterate through it and chech if we have responded to it before. by comparing the 3 first parameters
	//if not, calc priority_value and add an order response.
	//if yes, skip

	//Deleting responses that are not received_order_requests with the same sender id
	//Iterate over each order in order_response
	/*
		for j := 0; j < len(order_response); j++ {
			existing_order_response := order_response[j]
			// Check if the order is in received_order_requests
			found := false
			for _, received := range received_order_requests {
				if existing_order_response.Original_elevator == received.Original_elevator && existing_order_response.Time == received.Time {
					found = true
					break
				}
			}
			// If the order is not in received_order_requests, remove it from order_response
			if !found {
				//Dette er sånn man fjærner element j fra en liste
				fmt.Println("Removing a old response")
				order_response = append(order_response[:j], order_response[j+1:]...)
				j-- // Decrement j to account for the removed element
			}
		}*/
	//NEW
	var filtered_order_response []OrderRequests

	// Iterate over each order in order_response
	for _, existing_order_response := range order_response {
		found := false

		for _, received := range received_order_requests {
			if existing_order_response.Original_elevator == received.Original_elevator && existing_order_response.Time == received.Time && existing_order_response.Original_elevator == received_from_id {
				found = true
				break
			}
		}
		// Add the order only if it's found in received_order_requests
		if found {
			filtered_order_response = append(filtered_order_response, existing_order_response)
		} else {
			fmt.Println("Removing an old response")
		}
	}

	// Update order_response with the filtered slice
	order_response = filtered_order_response

	//Generating new response
	for _, receved_order := range received_order_requests {
		//for each receved_order
		request_has_a_response := false
		for _, order_resp := range order_response {
			//is that a receved order request has a respons in order_response go next
			request_has_a_response = false
			if order_resp.Original_elevator == receved_order.Original_elevator && order_resp.Time == receved_order.Time {
				// The order has already been responded to, so skip it
				request_has_a_response = true
				break
			}
		}
		// Calculate priority_value and add an order response
		if !request_has_a_response {
			fmt.Println("Generating new response")
			priority_value := fsm.Calculate_priority_value(receved_order.Button_event) // IKKE VELDIG PENT AT DENNE BLIR BRUKT HER!
			receved_order.Responder_elevator = ID
			receved_order.Elevator_priority_value_respond = priority_value
			order_response = append(order_response, receved_order)
			// ...
		}

	}

}

func Receved_order_response(received_order_responses []OrderRequests) {
	//fmt.Println("Receved_order_response, id: ")
	for _, order := range received_order_responses {
		if order.Original_elevator == ID {
			// Found a order respond that responds to one of our requests
			//Check if there still is possible to submit repons, by cheching if the chennel is down
			respnd_chan, exist := my_request_responses_chans[order.Time]
			if exist {
				fmt.Println("Send response to wait go, from id: ", order.Responder_elevator)
				response := SingleResponse{order.Responder_elevator, order.Elevator_priority_value_respond}
				respnd_chan <- response
			} else {
				fmt.Println("Response no longer valid")
			}
		}
	}
}
