package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"os"
	"strings"
)

func GetMotorDirectionFromDirn(dirn Dirn) elevio.MotorDirection {
	switch dirn {
	case DirUp:
		return elevio.MotorUp
	case DirDown:
		return elevio.MotorDown
	case DirStop:
		return elevio.MotorStop
	default:
		return elevio.MotorStop
	}
}

type Dirn int

const (
	DirUp   Dirn = 1
	DirDown Dirn = -1
	DirStop Dirn = 0
)

type ElevatorBehaviour int

const (
	ElevIdle     ElevatorBehaviour = 0
	ElevDoorOpen ElevatorBehaviour = 1
	ElevMoving   ElevatorBehaviour = 2
)

type Elevator struct {
	ID                 int
	Floor              int
	Dirn               Dirn
	Behaviour          ElevatorBehaviour
	Requests           [NumFloors][NumButtons]bool
	DoorOpenDuration_s float64
	Obstruction        bool
	MotorError         bool
}

// Elevator initializer function
func NewElevator() Elevator {
	var elevatorSetup Elevator = Elevator{
		ID:                 0,
		Floor:              -1,                            // Uninitialized floor
		Dirn:               DirStop,                       // Not moving
		Behaviour:          ElevIdle,                      // Idle state
		Requests:           [NumFloors][NumButtons]bool{}, // No requests initially
		DoorOpenDuration_s: 3.0,
		Obstruction:        false,
		MotorError:         false,
	}

	return elevatorSetup
}

func PrintElevator(elevator Elevator) {
	fmt.Printf("\n\nElevator:\n")
	fmt.Printf("ID: %d\n", elevator.ID)
	fmt.Printf("Floor: %d\n", elevator.Floor)
	fmt.Printf("Direction: %v\n", elevator.Dirn)
	fmt.Printf("Behaviour: %v\n", elevator.Behaviour)
	fmt.Printf("Requests: %v\n", elevator.Requests)
	fmt.Printf("Obstruction: %v\n", elevator.Obstruction)
}

func SetObstructionStatus(status bool) {
	setElevatorObtruction(status)
	if status {
		StartObstrTimer()
	} else {
		StopObstrTimer()
	}
	fmt.Println("Obstruction status set to:", status)
}

// Fucntion that sets all requests and lights to false
func ClearAllRequests() {
	elevator.Requests = [NumFloors][NumButtons]bool{}
	setAllLights(elevator)
}

// Function to store the last column of requests (cab calls) in a file, replacing the first row every time
func storeCabRequests(elevator Elevator) {
	file, err := os.OpenFile("./fsm/cab_requests.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("(Store) Error opening file:", err)
		return
	}
	defer file.Close()

	// Collect cab request states (last column of requests)
	var cabRequests []string
	for floor := 0; floor < NumFloors; floor++ {
		if elevator.Requests[floor][NumButtons-1] {
			cabRequests = append(cabRequests, "1")
		} else {
			cabRequests = append(cabRequests, "0")
		}
	}

	// Write the single row (overwrite file every time)
	_, err = file.WriteString(strings.Join(cabRequests, ",") + "\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// Function to read cab requests from a file and update the elevator's request state
func loadCabRequests(elevator *Elevator) {
	// Open the cab requests file
	file, err := os.Open("./fsm/cab_requests.txt")
	if err != nil {
		fmt.Println("(Load) Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the entire file content
	var content string
	_, err = fmt.Fscan(file, &content)
	if err != nil {
		fmt.Println("(Load) Error reading file:", err)
		return
	}

	// Split the content by commas to get individual button states
	cabRequests := strings.Split(content, ",")
	if len(cabRequests) != NumFloors {
		fmt.Println("(Load) Invalid file format: number of entries doesn't match the number of floors")
		return
	}

	// Update the elevator's request state based on the file content
	for floor := 0; floor < NumFloors; floor++ {
		if cabRequests[floor] == "1" {
			elevator.Requests[floor][NumButtons-1] = true
		} else {
			elevator.Requests[floor][NumButtons-1] = false
		}
	}

	fmt.Println("Cab requests loaded successfully")
}
