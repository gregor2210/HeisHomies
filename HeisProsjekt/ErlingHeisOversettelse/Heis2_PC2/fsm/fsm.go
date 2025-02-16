package fsm

import (
	"Driver-go/elevio"
	"fmt"
)

var elevator Elevator = NewElevator()

//ElevatorOUtputDevice er den utdelte go driverern!

func setAllLights(elevator Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, elevator.requests[floor][btn])
		}
	}
}

// func fsm_onInitBetweenFloors() {
// 	elevio.SetMotorDirection(elevio.MD_Down)
// 	elevator.dirn = elevio.MD_Down
// 	elevator.behaviour = EB_Moving
// }

func Fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	fmt.Printf("\n\nfsm_onRequestButtonPress(%d)\n", btn_floor)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) { // Hvis heisen allerede er i etasjen og knappen er trykket inn
			TimerStart(elevator.doorOpenDuration_s) // Start dørtimeren på nytt
		} else {
			elevator.requests[btn_floor][btn_type] = true // Ellers legg forespørselen til i køen
		}
	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = true // I bevegelse, så ikke gjøre annet enn å legge til i køen

	case EB_Idle:
		elevator.requests[btn_floor][btn_type] = true                   // Heisen er i ro, må finne retning og starte bevegelse
		var pair DirnBehaviourPair = requests_chooseDirection(elevator) // Velg retning basert på forespørsler
		elevator.dirn = pair.dirn                                       // Oppdater retning
		elevator.behaviour = pair.behaviour                             // Oppdater tilstand
		switch pair.behaviour {
		case EB_DoorOpen: // Hvis heisen skal stoppe i etasjen den står i
			elevio.SetDoorOpenLamp(true)
			TimerStart(elevator.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)

		case EB_Moving: // Hvis heisen skal starte bevegelse
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.dirn)) // Start motoren

		case EB_Idle:
		}
	}

	setAllLights(elevator) // Oppdater lysindikatorene

	fmt.Println("\nNew state:")
}

func Fsm_onFloorArrival(newFloor int) {
	fmt.Printf("\n\nfsm_onFloorArrival(%d)\n", newFloor)
	elevator.floor = newFloor

	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) { // Hvis heisen skal stoppe i etasjen den er i, enten fordi request i riktig retning i etasjen eller cab, eller ingen flere forespørsler
			elevio.SetMotorDirection(elevio.MD_Stop)          // Stopp motoren
			elevio.SetDoorOpenLamp(true)                      // Åpne døren
			elevator = requests_clearAtCurrentFloor(elevator) // Rydd opp forespørslene i etasjen
			TimerStart(elevator.doorOpenDuration_s)
			setAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}
	default:
	}

	fmt.Println("\nNew state:")
}

func Fsm_onDoorTimeout() {
	fmt.Printf("\n\nfsm_onDoorTimeout()\n")

	switch elevator.behaviour {
	case EB_DoorOpen:
		// Velg retning basert på forespørsler
		var pair DirnBehaviourPair = requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			// Restart dørtimeren
			TimerStart(elevator.doorOpenDuration_s)

			// Rydd opp forespørslene i nåværende etasje
			elevator = requests_clearAtCurrentFloor(elevator)

			// Oppdater alle lysindikatorer
			setAllLights(elevator)

		case EB_Moving:
			// Start motoren
			elevio.SetDoorOpenLamp(false) // Slå av dørindikatoren
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.dirn))
			fmt.Printf("Motor started moving in direction: %v\n", elevator.dirn)

		case EB_Idle:
			// Ingen forespørsler igjen, sett til idle
			elevio.SetDoorOpenLamp(false)
		}

	default:
		// Ingenting å gjøre hvis tilstanden ikke er EB_DoorOpen
	}

	fmt.Println("\nNew state:")
}
