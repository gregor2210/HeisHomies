package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"fmt"
)

type WorldviewPackage struct {
	ElevatorID     int
	Elevator       fsm.Elevator //refrence
	Order_requeset []OrderRequests
	Order_response []OrderRequests
	Order          elevio.ButtonEvent
	OrderBool      bool
}

func NewWorldviewPackage(ElevatorID int, elevator_ fsm.Elevator) WorldviewPackage {
	pending_orders := Get_pending_orders()
	responses := Get_order_respons()
	return WorldviewPackage{
		ElevatorID:     ElevatorID,
		Elevator:       elevator_,
		Order_requeset: pending_orders,
		Order_response: responses,
		OrderBool:      false,
	}
}

func PrintWorldview(worldView WorldviewPackage) {
	fmt.Print("WorldView")
	fmt.Println("Elevator ID: ", worldView.ElevatorID)

}
