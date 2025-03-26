package connectivity

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"sync"
)

var (
	// worldViewBackup[i] is the last known state of elevator i
	worldViewBackup      [NumElevators]WorldviewPackage
	worldViewBackupMutex sync.Mutex
)

func StoreWorldview(id int, worldview WorldviewPackage) {
	worldViewBackupMutex.Lock()
	defer worldViewBackupMutex.Unlock()
	worldViewBackup[id] = worldview

}

func GetWorldView(id int) WorldviewPackage {
	worldViewBackupMutex.Lock()
	defer worldViewBackupMutex.Unlock()
	return worldViewBackup[id]
}

// Returns true if the given order exists in any online elevator, including self
func DoesOrderExist(buttonEvent elevio.ButtonEvent) bool {
	floor := buttonEvent.Floor
	var button int = int(buttonEvent.Button) // 0 hallup, 1 halldown
	onlineElevatorIDs := GetAllOnlineIds()

	// Iterate through all online elevators, including itself, to check for button events
	for _, id := range onlineElevatorIDs {
		if id == ID {
			requests := fsm.GetElevatorStruct().Requests
			if requests[floor][button] {
				return true
			}

		} else {
			worldView := GetWorldView(id)
			requests := worldView.Elevator.Requests

			if requests[floor][button] {
				return true
			}
		}
	}

	return false
}
