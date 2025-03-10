package connectivity

import (
	"flag"
	"log"
	"net/http"
)

const (
	NR_OF_ELEVATORS = 3
	// Timeout for receiving UDP messages
	TIMEOUT = 3

	// Worldview max package size
	PACKAGE_SIZE = 1500
)

var (
	ID int //deafult 0
)

func init() {
	port_debug := flag.String("pprof-port", "6060", "Port for pprof")

	go func() {
		log.Println(http.ListenAndServe("localhost:"+*port_debug, nil))
	}()

	flag.IntVar(&ID, "id", 0, "Specify the id with -id")
	flag.Parse()

}
