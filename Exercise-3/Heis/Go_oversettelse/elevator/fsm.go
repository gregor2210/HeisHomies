package elevator

import (
	"fmt"
	"project/elevator"
	"project/elevator/requests"
	"project/elevator/timer"
)

// Elevator-tilstand og OutputDevice
var (
	elevator      elevator.Elevator
	outputDevice  elevator.ElevOutputDevice
)

// Initialiserer FSM
func init() {
	elevator = elevator.UninitializedElevator()

	// Henter konfigurasjon fra fil
	conLoad("elevator.con", map[string]interface{}{
		"doorOpenDuration_s":      &elevator.Config.DoorOpenDuration,
		"clearRequestVariant":     &elevator.Config.ClearRequestVariant,
	})

	// Initialiserer OutputDevice
	outputDevice = elevator.GetOutputDevice()
}

// Setter alle lysene basert på heisens tilstand
func setAllLights(e elevator.Elevator) {
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			outputDevice.RequestButtonLight(floor, btn, e.Requests[floor][btn])
		}
	}
}

// Kjør når heisen starter mellom etasjer
func OnInitBetweenFloors() {
	outputDevice.MotorDirection(elevator.D_Down)
	elevator.Dirn = elevator.D_Down
	elevator.Behaviour = elevator.EB_Moving
}

// Kjør når en knapp trykkes
func OnRequestButtonPress(btnFloor int, btnType elevator.Button) {
	fmt.Printf("\n\nOnRequestButtonPress(%d, %s)\n", btnFloor, elevator.ButtonToString(btnType))
	elevator.PrintState()

	switch elevator.Behaviour {

	case elevator.EB_DoorOpen:
		if requests.ShouldClearImmediately(elevator, btnFloor, btnType) {
			timer.Start(elevator.Config.DoorOpenDuration)
		} else {
			elevator.Requests[btnFloor][btnType] = 1
		}

	case elevator.EB_Moving:
		elevator.Requests[btnFloor][btnType] = 1

	case elevator.EB_Idle:
		elevator.Requests[btnFloor][btnType] = 1
		pair := requests.ChooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			outputDevice.DoorLight(1)
			timer.Start(elevator.Config.DoorOpenDuration)
			elevator = requests.ClearAtCurrentFloor(elevator)

		case elevator.EB_Moving:
			outputDevice.MotorDirection(elevator.Dirn)

		case elevator.EB_Idle:
			// Ingen handling nødvendig
		}
	}

	setAllLights(elevator)

	fmt.Println("\nNew state:")
	elevator.PrintState()
}

// Kjør når heisen ankommer en ny etasje
func OnFloorArrival(newFloor int) {
	fmt.Printf("\n\nOnFloorArrival(%d)\n", newFloor)
	elevator.PrintState()

	elevator.Floor = newFloor
	outputDevice.FloorIndicator(newFloor)

	switch elevator.Behaviour {

	case elevator.EB_Moving:
		if requests.ShouldStop(elevator) {
			outputDevice.MotorDirection(elevator.D_Stop)
			outputDevice.DoorLight(1)
			elevator = requests.ClearAtCurrentFloor(elevator)
			timer.Start(elevator.Config.DoorOpenDuration)
			setAllLights(elevator)
			elevator.Behaviour = elevator.EB_DoorOpen
		}

	default:
		// Ingen handling nødvendig
	}

	fmt.Println("\nNew state:")
	elevator.PrintState()
}

// Kjør når dørtidsavbrudd skjer
func OnDoorTimeout() {
	fmt.Printf("\n\nOnDoorTimeout()\n")
	elevator.PrintState()

	switch elevator.Behaviour {

	case elevator.EB_DoorOpen:
		pair := requests.ChooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case elevator.EB_DoorOpen:
			timer.Start(elevator.Config.DoorOpenDuration)
			elevator = requests.ClearAtCurrentFloor(elevator)
			setAllLights(elevator)

		case elevator.EB_Moving, elevator.EB_Idle:
			outputDevice.DoorLight(0)
			outputDevice.MotorDirection(elevator.Dirn)
		}
	}

	fmt.Println("\nNew state:")
	elevator.PrintState()
}