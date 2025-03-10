package connectivity

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

const (
	NR_OF_ELEVATORS = 3

	// Timeout for receiving and sending messages. If ether exside we consider connection lost.
	TIMEOUT = 3 //seconds

	// Worldview max package size
	PACKAGE_SIZE = 1500
)

var (
	ID int //deafult 0
)

func init() {
	//Setting up the flags for the code
	flag.IntVar(&ID, "id", 0, "Specify the id with -id")
	//FOr debugging purpoeses
	go func() {
		log.Println(http.ListenAndServe("localhost:"+strconv.Itoa(6060+ID), nil))
	}()

	flag.Parse()

}
