package connectivity

import "flag"

const (
	NumElevators = 3
	// TimeOut for receiving UDP messages
	TimeOut = 3

	// Worldview max package size
	MaxPacketSize = 1500
)

var (
	ID int //deafult 0
)

func init() {
	flag.IntVar(&ID, "id", 0, "Specify the id with -id")
	flag.Parse()
}
