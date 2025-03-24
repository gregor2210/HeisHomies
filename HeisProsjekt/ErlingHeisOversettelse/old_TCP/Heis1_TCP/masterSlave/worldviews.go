package masterSlave

import "log"

var (
	worldView_elv_1 WorldviewPackage
	worldView_elv_2 WorldviewPackage
	worldView_elv_3 WorldviewPackage
)

func StoreWorldview(ID int, worldview WorldviewPackage) {
	switch ID {
	case 1:
		worldView_elv_1 = worldview
	case 2:
		worldView_elv_2 = worldview
	case 3:
		worldView_elv_3 = worldview
	}

}

func GetWorldView(ID int) WorldviewPackage {
	switch ID {
	case 1:
		return worldView_elv_1
	case 2:
		return worldView_elv_2
	case 3:
		return worldView_elv_3
	}
	log.Fatalf("Error GetWorldView: Invalid ID%v", ID)
	return WorldviewPackage{}
}
