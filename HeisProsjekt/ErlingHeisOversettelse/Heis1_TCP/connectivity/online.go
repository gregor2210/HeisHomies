package connectivity

import (
	"fmt"
	"log"
	"sync"
)

var (
	// isOnline[i] is true if elevator i has recently sent a message
	isOnline       = [NR_OF_ELEVATORS]bool{}
	isOnline_mutex sync.Mutex

	// Channel to notfiy if an elevator goes offline
	offline_update_chan chan int
)

func init() {
	isOnline[ID] = true
	offline_update_chan = make(chan int)
}

func Online_setup(offline_update_chan_ chan int) {
	offline_update_chan = offline_update_chan_
}

func SetElevatorOnline(elevatorID int) {
	isOnline_mutex.Lock()
	defer isOnline_mutex.Unlock()

	if elevatorID >= 0 && elevatorID < NR_OF_ELEVATORS {

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
	fmt.Println("Returning from SetElevatorOnline")
}

func SetElevatorOffline(elevatorID int) {
	isOnline_mutex.Lock()
	defer isOnline_mutex.Unlock()

	if elevatorID >= 0 && elevatorID < len(isOnline) {
		// Only print if the elevator was previously online and is now set to offline
		if isOnline[elevatorID] {
			isOnline[elevatorID] = false
			offline_update_chan <- elevatorID

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
	fmt.Println("Returning from SetElevatorOffline")
}

// Return true if elevator is online
func IsOnline(elevatorID int) bool {
	isOnline_mutex.Lock()
	defer isOnline_mutex.Unlock()
	if elevatorID >= 0 && elevatorID < len(isOnline) {
		return isOnline[elevatorID]
	}

	log.Fatal("Not valid elevatorID when IsOnline()", elevatorID)
	return false
}

func PrintIsOnline() {

	isOnline_mutex.Lock()
	defer isOnline_mutex.Unlock()
	for i, online := range isOnline {
		status := "offline"
		if online {
			status = "online"
		}
		log.Printf("Elevator %d is %s\n", i, status)
	}
}

func Self_only_online() bool {
	isOnline_mutex.Lock()
	defer isOnline_mutex.Unlock()
	for i := 0; i < NR_OF_ELEVATORS; i++ {
		if i != ID {
			if isOnline[i] {
				return false
			}
		}
	}
	return true
}

func Get_all_online_ids() []int {
	isOnline_mutex.Lock()
	defer isOnline_mutex.Unlock()

	var online_elevators []int
	for i, online := range isOnline {
		if online {
			online_elevators = append(online_elevators, i)
		}
	}
	return online_elevators
}
