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

// Calculating the priority value, return the value
func CalculatePriorityValue(buttonEvent elevio.ButtonEvent, e fsm.Elevator) int {
	// This function calculates the priority value.
	// This value indicates how much an elevator wants a request; a higher value means it wants it more.
	// It starts with a max value and then proceeds to subtract values.
	// It will subtract 10 for every floor it is away from the button press.
	//
	// It will consider the worst-case number of floors to reach the button.
	// For example, if the elevator is on the 3rd floor and moving down to the 2nd floor, and an event occurs at the 4th floor,
	// it will assume the worst-case scenario where it travels to the bottom floor before going up to the button.
	// Therefore, it will subtract 10 for each floor on its way down and back up to the button.

	// Larger value == higher order priority.
	requestFloor := buttonEvent.Floor

	NumFloorsMinus1 := fsm.NumFloors - 1
	//Calculate how much this elevator wants this request.
	priorityValue := 2 * 10 * NumFloorsMinus1 // max value

	//DÅRLIG VERSJON, TAR ABSOLUTT AVSTAND
	deltaFloor := requestFloor - e.Floor

	subVal := int(math.Abs(float64(deltaFloor))) * 10
	//if elevator dosen ot have a moving dirn
	if int(e.Dirn) == 0 {
		// Elevator ha no order
		//-

	} else if deltaFloor < 0 && int(e.Dirn) < 0 {
		// Elevator moves down toward the button, down

		if buttonEvent.Button == elevio.BtnHallDown {
			// Button moves in the same direction as the elevator: down
			//-

		} else {
			//  Button points in the opposite direction of the elevator: up
			subVal += buttonEvent.Floor * 2 * 10

		}

	} else if deltaFloor > 0 && int(e.Dirn) < 0 {
		// Elevator moves down, away from the button
		subVal += 2 * e.Floor * 10

		if buttonEvent.Button == elevio.BtnHallUp {
			// Button points up
			//-
		} else {
			// Button points down
			subVal += (NumFloorsMinus1 - buttonEvent.Floor) * 2 * 10

		}

	} else if deltaFloor > 0 && int(e.Dirn) > 0 {
		// Elevator moves up toward the button, up

		if buttonEvent.Button == elevio.BtnHallUp {
			// Button points in the same direction: up
			//-

		} else {
			// Knapp peker motsatt vei NED
			subVal += (NumFloorsMinus1 - buttonEvent.Floor) * 2 * 10
		}
	} else if deltaFloor < 0 && int(e.Dirn) > 0 {
		// Elevator moves up, away from the button
		subVal += (NumFloorsMinus1 - e.Floor) * 2 * 10

		if buttonEvent.Button == elevio.BtnHallUp {
			// Button points up
			subVal += buttonEvent.Floor * 2 * 10
		} else {
			// Button points down
			//-
		}
	}

	priorityValue -= subVal

	//fmt.Println("Priority value:", priorityValue)
	return priorityValue
}

func NewOrder(buttonEvent elevio.ButtonEvent) {
	// Figure out who should take which order.
	// Sends the order to the selected elevator

	if DoesOrderExist(buttonEvent) {
		fmt.Println("Order allready exist")
		return
	}

	var priorityValueIDIndex [NumElevators]int
	var priorityValueToSort []int
	onlineElevatorID := GetAllOnlineIds()
	for _, id := range onlineElevatorID {
		//Gets list of all priority values
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
		// Find the elevator ID that will receive the order
		idOfElevatorThatWillGetOrder := ID
		for i, v := range priorityValueIDIndex {
			if v == priorityValue {
				idOfElevatorThatWillGetOrder = i
				break
			}
		}

		if idOfElevatorThatWillGetOrder == ID {
			//Send reqeust to self
			fsm.FsmOnRequestButtonPress(buttonEvent.Floor, buttonEvent.Button) // Ikke så fint at dnne er her
			didOrderGetSent = true
			break

		} else if SendOrderToSpecificElevator(idOfElevatorThatWillGetOrder, buttonEvent) {
			//Try to send order to elevator with id
			fmt.Println("ID: ", idOfElevatorThatWillGetOrder, "Got the order!")
			didOrderGetSent = true
			// Update its wv backup with this order
			setOrderOnbackup(idOfElevatorThatWillGetOrder, buttonEvent)
			break
		}

	}
	if !didOrderGetSent {
		fmt.Println("No elevator could take the order")
	}

}
