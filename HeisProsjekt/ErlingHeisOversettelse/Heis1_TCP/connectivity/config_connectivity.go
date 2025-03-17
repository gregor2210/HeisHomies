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

	// USE_IPS is set to true if you are gonna use different computer.
	// Remember to set correct ips in IPs
	USE_IPS = false
)

var (
	ID int //deafult 0

	//Ips of the different computers.
	// computer with id 0 has inex 0 in this list
	IPs = [NR_OF_ELEVATORS]string{"10.100.23.82", "10.100.23.32", "10.100.23.29"}
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
