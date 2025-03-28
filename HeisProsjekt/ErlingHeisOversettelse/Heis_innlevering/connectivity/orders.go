package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
	"math"
	"sort"
)

type OrderRequests struct {
	OriginalElevator             int
	ResponderElevator            int
	UniqueID                     string
	ButtonEvent                  elevio.ButtonEvent
	ElevatorPriorityValueRespond int
}

func PrintOrderRequest(orders []OrderRequests) {
	fmt.Println("Pending Orders:")
	for _, order := range orders {
		fmt.Printf("Original Elevator: %d, Responder Elevator: %d, Time: %s, Priority: %d\n",
			order.OriginalElevator, order.ResponderElevator, order.UniqueID, order.ElevatorPriorityValueRespond)
	}
}

func CalculatePriorityValue(buttonEvent elevio.ButtonEvent, e fsm.Elevator) int {
	// Calculates how much the elevator wants the request:
	// 1. Starts at max value.
	// 2. Subtracts 10 per floor from request.
	// 3. Penalizes detours and opposite direction, assuming worst-case travel.

	requestFloor := buttonEvent.Floor

	NumFloorsMinus1 := fsm.NumFloors - 1
	priorityValue := 2 * 10 * NumFloorsMinus1 // max value
	deltaFloor := requestFloor - e.Floor
	subVal := int(math.Abs(float64(deltaFloor))) * 10

	// Elevator ha no order
	if int(e.Dirn) == 0 {

		// Moving down toward request
	} else if deltaFloor < 0 && int(e.Dirn) < 0 {

		// Same direction: down
		if buttonEvent.Button == elevio.BtnHallDown {

		} else {
			// Opposite direction: up
			subVal += buttonEvent.Floor * 2 * 10

		}

		// Moving down, away from request
	} else if deltaFloor > 0 && int(e.Dirn) < 0 {
		subVal += 2 * e.Floor * 10

		// Request is up
		if buttonEvent.Button == elevio.BtnHallUp {

		} else {
			// Request is down
			subVal += (NumFloorsMinus1 - buttonEvent.Floor) * 2 * 10

		}
		// Moving up toward request
	} else if deltaFloor > 0 && int(e.Dirn) > 0 {

		// Same direction: up
		if buttonEvent.Button == elevio.BtnHallUp {

		} else {
			// Opposite direction: down
			subVal += (NumFloorsMinus1 - buttonEvent.Floor) * 2 * 10
		}
		// Moving up, away from request
	} else if deltaFloor < 0 && int(e.Dirn) > 0 {
		subVal += (NumFloorsMinus1 - e.Floor) * 2 * 10

		// Request is up
		if buttonEvent.Button == elevio.BtnHallUp {
			subVal += buttonEvent.Floor * 2 * 10
		} else {
			// Request is down
		}
	}

	priorityValue -= subVal
	return priorityValue
}

func NewOrder(buttonEvent elevio.ButtonEvent) {
	// Figure out who should take which order
	// Sends the order to the selected elevator

	if DoesOrderExist(buttonEvent) {
		fmt.Println("Order allready exist")
		return
	}

	// Get the list of all priority values
	var priorityValueIDIndex [NumElevators]int
	var priorityValueToSort []int
	onlineElevatorID := GetAllOnlineIds()
	for _, id := range onlineElevatorID {
		var elevator fsm.Elevator
		if id == ID {
			elevator = fsm.GetElevatorStruct()

		} else {
			elevator = GetWorldView(id).Elevator

		}

		priorityValue := CalculatePriorityValue(buttonEvent, elevator)
		priorityValueIDIndex[id] = priorityValue
		priorityValueToSort = append(priorityValueToSort, priorityValue)
	}

	// Sorting priorityValueToSort in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(priorityValueToSort)))

	didOrderGetSent := false
	// Finding the elevator ID with the highest priority value and attempting to send the order to that elevator
	for _, priorityValue := range priorityValueToSort {
		idOfElevatorThatWillGetOrder := ID
		for i, v := range priorityValueIDIndex {
			if v == priorityValue {
				idOfElevatorThatWillGetOrder = i
				break
			}
		}
		//Send reqeust to self
		if idOfElevatorThatWillGetOrder == ID {
			fsm.FsmOnRequestButtonPress(buttonEvent.Floor, buttonEvent.Button) // Ikke så fint at dnne er her
			didOrderGetSent = true
			break

			//Try sending order to selected elevator
		} else if SendOrderToSpecificElevator(idOfElevatorThatWillGetOrder, buttonEvent) {

			fmt.Println("ID: ", idOfElevatorThatWillGetOrder, "Got the order!")
			didOrderGetSent = true

			// Update its worldview backup with this order
			//setOrderOnbackup(idOfElevatorThatWillGetOrder, buttonEvent)
			break
		}

	}
	if !didOrderGetSent {
		fmt.Println("No elevator could take the order")
	}

}
