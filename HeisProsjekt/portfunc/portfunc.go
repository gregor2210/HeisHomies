package makeport

import (
	"fmt"
)

const (
	ID              = 0
	Nr_of_elevators = 3
)

var (
	UDP_world_view_send_port    []string
	UDP_other_pc_ips            []string
	UDP_world_view_receive_port []string
)

func makeports() {
	var elevatorIDs []int
	for i := 0; i < Nr_of_elevators; i++ {
		if i != ID {
			elevatorIDs = append(elevatorIDs, i)
		}
	}

	for _, elm := range elevatorIDs {
		UDP_world_view_send_port = append(UDP_world_view_send_port, fmt.Sprintf("80%d%d", ID, i))
		UDP_world_view_receive_port = append(UDP_world_view_receive_port, fmt.Sprintf("80%d%d", i, ID))
	}
}
