package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"sync"
)

var (
	elevator      Elevator = NewElevator()
	elevatorMutex sync.Mutex
)

func GetElevatorStruct() Elevator {
	elevatorMutex.Lock()
	defer elevatorMutex.Unlock()
	return elevator
}

//ElevatorOUtputDevice er den utdelte go driverern!

func setAllLights(elevator Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, elevator.Requests[floor][btn])
		}
	}
}

// func fsm_onInitBetweenFloors() {
// 	elevio.SetMotorDirection(elevio.MotorDown)
// 	elevator.dirn = elevio.MotorDown
// 	elevator.behaviour = ElevMoving
// }

func FsmOnRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Printf("\n\nFsmOnRequestButtonPress(%d)\n", btnFloor)

	switch elevator.Behaviour {
	case ElevDoorOpen:
		if requestsShouldClearImmediately(elevator, btnFloor, btnType) { // Hvis heisen allerede er i etasjen og knappen er trykket inn
			TimerStart(elevator.DoorOpenDuration_s) // Start dørtimeren på nytt
		} else {
			elevator.Requests[btnFloor][btnType] = true // Ellers legg forespørselen til i køen
		}
	case ElevMoving:
		elevator.Requests[btnFloor][btnType] = true // I bevegelse, så ikke gjøre annet enn å legge til i køen

	case ElevIdle:
		elevator.Requests[btnFloor][btnType] = true                    // Heisen er i ro, må finne retning og starte bevegelse
		var pair DirnBehaviourPair = requestsChooseDirection(elevator) // Velg retning basert på forespørsler
		elevator.Dirn = pair.dirn                                      // Oppdater retning
		elevator.Behaviour = pair.behaviour                            // Oppdater tilstand
		switch pair.behaviour {
		case ElevDoorOpen: // Hvis heisen skal stoppe i etasjen den står i
			elevio.SetDoorOpenLamp(true)
			TimerStart(elevator.DoorOpenDuration_s)
			elevator = requestsClearAtCurrentFloor(elevator)

		case ElevMoving: // Hvis heisen skal starte bevegelse
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn)) // Start motoren

		case ElevIdle:
		}
	}

	setAllLights(elevator) // Oppdater lysindikatorene

	fmt.Println("\nNew state:")
}

func FsmOnFloorArrival(newFloor int) {
	fmt.Printf("\n\nFsmOnFloorArrival(%d)\n", newFloor)
	elevator.Floor = newFloor

	elevio.SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case ElevMoving:
		if requestsShouldStop(elevator) { // Hvis heisen skal stoppe i etasjen den er i, enten fordi request i riktig retning i etasjen eller cab, eller ingen flere forespørsler
			elevio.SetMotorDirection(elevio.MotorStop)       // Stopp motoren
			elevio.SetDoorOpenLamp(true)                     // Åpne døren
			elevator = requestsClearAtCurrentFloor(elevator) // Rydd opp forespørslene i etasjen
			TimerStart(elevator.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = ElevDoorOpen
		}
	default:
	}

	fmt.Println("\nNew state:")
}

func FsmOnDoorTimeOut() {
	fmt.Printf("\n\nFsmOnDoorTimeOut()\n")

	switch elevator.Behaviour {
	case ElevDoorOpen:
		// Velg retning basert på forespørsler
		var pair DirnBehaviourPair = requestsChooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch elevator.Behaviour {
		case ElevDoorOpen:
			// Restart dørtimeren
			TimerStart(elevator.DoorOpenDuration_s)

			// Rydd opp forespørslene i nåværende etasje
			elevator = requestsClearAtCurrentFloor(elevator)

			// Oppdater alle lysindikatorer
			setAllLights(elevator)

		case ElevMoving:
			// Start motoren
			elevio.SetDoorOpenLamp(false) // Slå av dørindikatoren
			elevio.SetMotorDirection(GetMotorDirectionFromDirn(elevator.Dirn))
			fmt.Printf("Motor started moving in direction: %v\n", elevator.Dirn)

		case ElevIdle:
			// Ingen forespørsler igjen, sett til idle
			elevio.SetDoorOpenLamp(false)
		}

	default:
		// Ingenting å gjøre hvis tilstanden ikke er ElevDoorOpen
	}

	fmt.Println("\nNew state:")
}
