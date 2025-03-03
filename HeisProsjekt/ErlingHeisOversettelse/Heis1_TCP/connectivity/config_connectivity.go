package connectivity

import "flag"

const (
	NR_OF_ELEVATORS = 4
)

var (
	ID int //deafult 0
)

func init() {
	flag.IntVar(&ID, "id", 0, "Specify the id with -id")
	flag.Parse()
}

