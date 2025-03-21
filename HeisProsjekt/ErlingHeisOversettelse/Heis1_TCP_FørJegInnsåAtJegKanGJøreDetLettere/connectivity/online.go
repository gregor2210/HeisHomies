package connectivity

import (
	"fmt"
	"log"
	"sync"
)

// This file contains the status of which elevators are are online and ofline based on this elevators view

var (
	// isOnline is a list of online elevators
	// The index of the list is the elevator ID
	// The value is true if the elevator is online, false if it is offline
	// Online or offline is based on if we receve message from it or not
	isOnline      = [NumElevators]bool{}
	isOnlineMutex sync.Mutex

	offlineUpdateChan chan int
)

func init() {
	// Initialize the isOnline list
	isOnline[ID] = true
}

func OnlineSetup(offlineUpdateChan_ chan int) {
	offlineUpdateChan = offlineUpdateChan_
}

func get_isOnline() [NumElevators]bool {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	return isOnline
}

func set_isOnline(id int, state bool) {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	isOnline[id] = state
}

// AddElevatorOnline sets the elevator ID to online in the isOnline list
func SetElevatorOnline(elevatorID int) {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()

	//If a valid id
	if elevatorID >= 0 && elevatorID < NumElevators {

		// If is only to make the print only appare if there is a chainge in state
		if isOnline[elevatorID] {
			fmt.Println("Setting ElevatorID:", elevatorID, "to ONLINE!")

		}

		isOnline[elevatorID] = true

	} else {
		log.Fatal("Not valid elevatorID when SetElevatorOnline()", elevatorID)
	}
}

// RemoveElevatorOnline sets the elevator ID to offline in the isOnline list
func SetElevatorOffline(elevatorID int) {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	if elevatorID >= 0 && elevatorID < len(isOnline) {

		// If is only to make the print only appare if there is a chainge in state
		if isOnline[elevatorID] {
			fmt.Println("Setting ElevatorID:", elevatorID, "to OFLINE!")
		}

		isOnline[elevatorID] = false
	} else {
		log.Fatal("Not valid elevatorID when SetElevatorOffline():", elevatorID)
	}
}

// Return true if is online
func IsOnline(elevatorID int) bool {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	if elevatorID >= 0 && elevatorID < len(isOnline) {
		return isOnline[elevatorID]
	}

	log.Fatal("Not valid elevatorID when IsOnline()", elevatorID)
	return false
}

// PrintIsOnline prints the status of all elevators
func PrintIsOnline() {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	for i, online := range isOnline {
		status := "offline"
		if online {
			status = "online"
		}
		log.Printf("Elevator %d is %s\n", i, status)
	}
}

func SelfOnlyOnline() bool {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	for i := 0; i < NumElevators; i++ {
		if i != ID {
			if isOnline[i] {
				return false
			}
		}
	}
	//Self is only online
	return true
}

func GetAllOnlineIds() []int {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()

	var onlineElevators []int
	for i, online := range isOnline {
		if online {
			onlineElevators = append(onlineElevators, i)
		}
	}
	return onlineElevators
}

func get_lowest_online_id() int {
	//returning the online elevator with lowest id
	return GetAllOnlineIds()[0]
}
