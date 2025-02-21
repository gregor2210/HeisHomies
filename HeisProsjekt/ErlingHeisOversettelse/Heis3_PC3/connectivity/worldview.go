package connectivity

import "Driver-go/fsm"

type Worldview_package struct {
	Elevator_ID int
	Elevator    fsm.Elevator //refrence
}

func New_Worldview_package(elevator_id int, elevator fsm.Elevator) Worldview_package {
	return Worldview_package{
		Elevator_ID: elevator_id,
		Elevator:    elevator,
	}
}
