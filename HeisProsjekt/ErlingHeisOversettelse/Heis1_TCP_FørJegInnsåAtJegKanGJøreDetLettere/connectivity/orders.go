package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type OrderRequests struct {
	OriginalElevator             int
	ResponderElevator            int
	UniqueID                     string
	buttonEvent                  elevio.ButtonEvent
	ElevatorPriorityValueRespond int
}

type SingleResponse struct {
	receved_from_id int
	priorityValue   int
}

type DoneProcessedOrder struct {
	Responses        [NumElevators]int
	Sorted_responses []int
	Order            elevio.ButtonEvent
}

var (
	pending_orders []OrderRequests
	order_response []OrderRequests

	UniqueID_counter     = 0
	UniqueID_counter_max = 9999999

	my_request_responses_chans = make(map[string]chan SingleResponse)

	pending_orders_mutex sync.Mutex
	order_response_mutex sync.Mutex
	UniqueID_mutex       sync.Mutex

	order_to_send_chan chan DoneProcessedOrder
)

func Order_setup(order_to_send_chan_ chan DoneProcessedOrder) {
	order_to_send_chan = order_to_send_chan_
}

func PrintOrderRequest(orders []OrderRequests) {
	fmt.Println("Pending Orders:")
	for _, order := range orders {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.OriginalElevator, order.ResponderElevator, order.UniqueID, order.ElevatorPriorityValueRespond)
	}
}

func get_UniqueID() string {
	UniqueID_mutex.Lock()
	defer UniqueID_mutex.Unlock()
	UniqueID_counter++
	if UniqueID_counter == UniqueID_counter_max {
		UniqueID_counter = 0
	}
	return strconv.Itoa(UniqueID_counter)
}

func Get_pending_orders() []OrderRequests {
	pending_orders_mutex.Lock()
	//fmt.Println("pending_orders_mutex locked")
	defer pending_orders_mutex.Unlock()
	//fmt.Println("pending_orders_mutex UNlocked")
	return pending_orders
}

func Get_order_respons() []OrderRequests {
	//fmt.Println("order_response_mutex locked")
	order_response_mutex.Lock()
	defer order_response_mutex.Unlock()
	//fmt.Println("order_response_mutex UNlocked")
	return order_response
}

func NewOrder(buttonEvent elevio.ButtonEvent, priorityValue int) {
	fmt.Println("New Order")
	pending_orders_mutex.Lock()
	defer pending_orders_mutex.Unlock()

	//Setting up the new order
	now := time.Now()
	now_str := now.Format("2006/01/02 15:04:05.999999999")
	UniqueID := now_str + " " + get_UniqueID()
	fmt.Println("now_str", UniqueID)
	NewOrder := OrderRequests{OriginalElevator: ID, UniqueID: UniqueID, buttonEvent: buttonEvent}
	pending_orders = append(pending_orders, NewOrder)

	//setting up listening to repsond go rutine
	response_chan := make(chan SingleResponse)
	//add it to dict to later be able to send SingleResponse to the corrrect go rutine based on the now_str
	my_request_responses_chans[UniqueID] = response_chan

	//priorityValue := CalculatePriorityValue()
	go Wait_for_order_responses(UniqueID, response_chan, priorityValue, GetAllOnlineIds(), buttonEvent)
	fmt.Println("New Order added and waiting response")
}

func Wait_for_order_responses(UniqueID string, singel_response_chan chan SingleResponse, this_elevators_priorityValue int, onlineElevators []int, order elevio.ButtonEvent) {
	fmt.Println("Started go rutine Wait_for_order_responses")
	var responses [NumElevators]int
	//responses_to_sort is a variabel that can be sort later, the use it to find the indexe/id of elevator to get the order
	//it makes it possible to use builtinn sort func
	var responses_to_sort []int
	response_counter := 0

	responses[ID] = this_elevators_priorityValue
	responses_to_sort = append(responses_to_sort, this_elevators_priorityValue)
	response_counter++

	//Timeoout will kick in after given miliseconds. The elevator will give the respons to the highest values.
	var TimeOut <-chan time.Time
	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop() // Ensure the ticker stops when the program exits
	TimeOut = ticker.C

loop:
	for {
		select {
		case single_response := <-singel_response_chan:
			if responses[single_response.receved_from_id] == 0 {
				fmt.Println("Got a new response, added, from", single_response.receved_from_id)
				responses[single_response.receved_from_id] = single_response.priorityValue
				responses_to_sort = append(responses_to_sort, single_response.priorityValue)
				response_counter++
			}
			if response_counter == len(onlineElevators) {
				fmt.Println("All responses receved")
				delete(my_request_responses_chans, UniqueID)
				break loop
			}

		case <-TimeOut:
			delete(my_request_responses_chans, UniqueID)
			break loop
		}

	}
	fmt.Println("wait for respond loop finished!")

	//Remoing the this wait order respons pending order from pending_order
	pending_orders_mutex.Lock()
	filtered_orders := pending_orders[:0] // Keeps the same underlying array
	for _, v := range pending_orders {
		if v.UniqueID != UniqueID {
			filtered_orders = append(filtered_orders, v)
		}
	}
	pending_orders = filtered_orders
	pending_orders_mutex.Unlock()

	//Choose who will get the order. find highest value and then the correspoinding index
	sort.Sort(sort.Reverse(sort.IntSlice(responses_to_sort))) // Sort descending
	done_processed_order := DoneProcessedOrder{responses, responses_to_sort, order}
	order_to_send_chan <- done_processed_order
	fmt.Println("'wait for response' exiting fucntioin")
}

func SendOrderToSpesificElevator(done_processed_order DoneProcessedOrder) bool {
	fmt.Println("Sending order to spesiffic pc")
	order := done_processed_order.Order
	responses := done_processed_order.Responses
	sorted_responses := done_processed_order.Sorted_responses

	return_bool := false

	for _, priorityValue := range sorted_responses {

		//find id or elevator that will get
		idOfElevatorThatWillGetOrder := ID
		for i, v := range responses {
			if v == priorityValue {
				idOfElevatorThatWillGetOrder = i
				break
			}
		}

		if idOfElevatorThatWillGetOrder == ID {
			//Send reqeust to self
			return_bool = false
			break

		} else if SendOrderToSpecificElevator(idOfElevatorThatWillGetOrder, order) {
			fmt.Println("ID: ", idOfElevatorThatWillGetOrder, "Got the order!")
			return_bool = true
			break
		}

	}
	return return_bool

}

func Receved_order_requests(received_order_requests []OrderRequests, received_from_id int) {
	//fmt.Println("Received order request----------------------------------")
	/*fmt.Println("Pending Orders:")
	for _, order := range pending_orders {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.OriginalElevator, order.ResponderElevator, order.Time, order.ElevatorPriorityValueRespond)
	}

	fmt.Println("\nOrder Response:")
	for _, order := range order_response {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.OriginalElevator, order.ResponderElevator, order.Time, order.ElevatorPriorityValueRespond)
	}
	fmt.Println("\nRECEVED Order Response:")
	for _, order := range received_order_requests {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.OriginalElevator, order.ResponderElevator, order.Time, order.ElevatorPriorityValueRespond)
	}
	*/
	order_response_mutex.Lock()
	defer order_response_mutex.Unlock()
	//Receves the incomming order_request_array
	//Removes any respons to a request that is not in request with the same id/signature.
	// Itterate through it and chech if we have responded to it before. by comparing the 3 first parameters
	//if not, calc priorityValue and add an order response.
	//if yes, skip

	//Deleting responses that are not received_order_requests with the same sender id
	var filtered_order_response []OrderRequests

	// Iterate over each order in order_response
	for _, existing_order_response := range order_response {
		found := false

		//hvis response er fra en annen pc enn den som vi sjekker orderen fra nå, kan vi IKKE slette. Det betyr ar den er da true
		if existing_order_response.OriginalElevator != received_from_id {
			found = true

			// Den eksisterende responsen er fra samme pc som den requesten vi sjekker nå
		} else {
			for _, received := range received_order_requests {
				//hvis det ikke finnes en request som har samme id og time som den eksisterende orderen, betyr det at den er løst og kan fjernes.
				// if'en er true hvis den fortsatt eksisterer
				if existing_order_response.OriginalElevator == received.OriginalElevator && existing_order_response.UniqueID == received.UniqueID {
					found = true
					break
				}
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
			if order_resp.OriginalElevator == receved_order.OriginalElevator && order_resp.UniqueID == receved_order.UniqueID {
				// The order has already been responded to, so skip it
				request_has_a_response = true
				break
			}
		}
		// Calculate priorityValue and add an order response
		if !request_has_a_response {
			fmt.Println("Generating new response")
			priorityValue := fsm.CalculatePriorityValue(receved_order.buttonEvent) // IKKE VELDIG PENT AT DENNE BLIR BRUKT HER!
			receved_order.ResponderElevator = ID
			receved_order.ElevatorPriorityValueRespond = priorityValue
			order_response = append(order_response, receved_order)
			// ...
		}

	}

}

func Receved_order_response(received_order_responses []OrderRequests) {
	//fmt.Println("Receved_order_response, id: ")
	for _, order := range received_order_responses {
		if order.OriginalElevator == ID {
			// Found a order respond that responds to one of our requests
			//Check if there still is possible to submit repons, by cheching if the chennel is down
			respnd_chan, exist := my_request_responses_chans[order.UniqueID]
			if exist {
				fmt.Println("Send response to wait go, from id: ", order.ResponderElevator)
				response := SingleResponse{order.ResponderElevator, order.ElevatorPriorityValueRespond}
				respnd_chan <- response
			} else {
				fmt.Println("Response no longer valid")
			}
		}
	}
}
