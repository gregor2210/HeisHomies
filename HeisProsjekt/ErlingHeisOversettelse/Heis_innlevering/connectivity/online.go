package connectivity

import (
	"fmt"
	"log"
	"sync"
)

var (
	// isOnline[i] is true if elevator i has recently sent a message
	isOnline      = [NumElevators]bool{}
	isOnlineMutex sync.Mutex

	// Channel to notfiy if an elevator goes offline
	offlineUpdateChan chan int
)

func init() {
	isOnline[ID] = true
	offlineUpdateChan = make(chan int, 10) // Buffer size 10
}

func OnlineSetup(offlineUpdateChan_ chan int) {
	offlineUpdateChan = offlineUpdateChan_
}

func SetSelfOnline() {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()

	isOnline[ID] = true
}

func SetSelfOffline() {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()

	isOnline[ID] = false
}

func IsSelfOnline() bool {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()
	return isOnline[ID]
}

func SetElevatorOnline(elevatorID int) {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()

	if elevatorID >= 0 && elevatorID < NumElevators {

		// Only print if the elevator was previously offline and is now set to online
		if !isOnline[elevatorID] {
			isOnline[elevatorID] = true

			fmt.Println("Setting ElevatorID:", elevatorID, "to ONLINE!")
			for i, online := range isOnline {
				status := "offline"
				if online {
					status = "online"
				}
				log.Printf("Elevator %d is %s\n", i, status)
			}
		}

	} else {
		log.Fatal("Not valid elevatorID when SetElevatorOnline()", elevatorID)
	}

}

func SetElevatorOffline(elevatorID int) {
	isOnlineMutex.Lock()
	defer isOnlineMutex.Unlock()

	if elevatorID >= 0 && elevatorID < len(isOnline) {
		// Only print if the elevator was previously online and is now set to offline
		if isOnline[elevatorID] {
			isOnline[elevatorID] = false
			offlineUpdateChan <- elevatorID

			fmt.Println("Setting ElevatorID:", elevatorID, "to OFFLINE!")
			for i, online := range isOnline {
				status := "offline"
				if online {
					status = "online"
				}
				log.Printf("Elevator %d is %s\n", i, status)
			}
		}

	} else {
		log.Fatal("Not valid elevatorID when SetElevatorOffline():", elevatorID)
	}
}

// Return true if elevator is online
func IsOnline(elevatorID int) bool {
	isOnlineMutex.Lock()

	defer isOnlineMutex.Unlock()
	if elevatorID >= 0 && elevatorID < len(isOnline) {
		return isOnline[elevatorID]
	}

	log.Fatal("Not valid elevatorID when IsOnline()", elevatorID)
	return false
}

// Print the online status of all elevators
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

// Return true if self is the only online elevator
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
