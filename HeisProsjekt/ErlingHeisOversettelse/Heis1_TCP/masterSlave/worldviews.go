package masterSlave

import "log"

var (
	world_view_elv_1 Worldview_package
	world_view_elv_2 Worldview_package
	world_view_elv_3 Worldview_package
)

func Store_worldview(ID int, worldview Worldview_package) {
	switch ID {
	case 1:
		world_view_elv_1 = worldview
	case 2:
		world_view_elv_2 = worldview
	case 3:
		world_view_elv_3 = worldview
	}

}

func Get_worldview(ID int) Worldview_package {
	switch ID {
	case 1:
		return world_view_elv_1
	case 2:
		return world_view_elv_2
	case 3:
		return world_view_elv_3
	}
	log.Fatalf("Error Get_worldview: Invalid ID%v", ID)
	return Worldview_package{}
}
