package connectivity

import (
	"sync"
)

var (
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
