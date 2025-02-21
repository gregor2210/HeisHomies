package connectivity

import (
	"fmt"
	"log"
)

// This file contains the status of which elevators are are online and ofline based on this elevators view

var (
	// isOnline is a list of online elevators
	// The index of the list is the elevator ID
	// The value is true if the elevator is online, false if it is offline
	// Online or offline is based on if we receve message from it or not
	isOnline = []bool{false, false, false}
)

func init() {
	// Initialize the isOnline list
	isOnline[ID] = true
}

// AddElevatorOnline sets the elevator ID to online in the isOnline list
func SetElevatorOnline(elevatorID int) {
	if elevatorID >= 0 && elevatorID < len(isOnline) {

		// If is only to make the print only appare if there is a chainge in state
		if !IsOnline(elevatorID) {
			fmt.Println("Setting ElevatorID:", elevatorID, "to ONLINE!")
		}

		isOnline[elevatorID] = true

	} else {
		log.Fatal("Not valid elevatorID when SetElevatorOnline()", elevatorID)
	}
}

// RemoveElevatorOnline sets the elevator ID to offline in the isOnline list
func SetElevatorOffline(elevatorID int) {
	if elevatorID >= 0 && elevatorID < len(isOnline) {

		// If is only to make the print only appare if there is a chainge in state
		if IsOnline(elevatorID) {
			fmt.Println("Setting ElevatorID:", elevatorID, "to OFLINE!")
		}

		isOnline[elevatorID] = false
	} else {
		log.Fatal("Not valid elevatorID when SetElevatorOffline():", elevatorID)
	}
}

// Return true if is online
func IsOnline(elevatorID int) bool {
	if elevatorID >= 0 && elevatorID < len(isOnline) {
		return isOnline[elevatorID]
	}

	log.Fatal("Not valid elevatorID when IsOnline()", elevatorID)
	return false
}

// PrintIsOnline prints the status of all elevators
func PrintIsOnline() {
	for i, online := range isOnline {
		status := "offline"
		if online {
			status = "online"
		}
		log.Printf("Elevator %d is %s\n", i, status)
	}
}
