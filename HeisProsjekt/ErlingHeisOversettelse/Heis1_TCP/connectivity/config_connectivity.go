package connectivity

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

const (
	NR_OF_ELEVATORS = 3

	// Timeout for send and receive. If exceeded, connection is lost.
	TIMEOUT = 3 //seconds

	// Worldview max package size
	PACKAGE_SIZE = 1500

	// USE_IPS is set to true if you are gonna use different computer.
	// Remember to set correct ips in IPs
	USE_IPS = false
)

var (
	ID int //default 0

	// IPs of all elevators
	IPs = [NR_OF_ELEVATORS]string{"10.100.23.28", "10.100.23.32", "10.100.23.29"}
)

func init() {
	// Set up program flags
	flag.IntVar(&ID, "id", 0, "Specify the id with -id")

	//For debugging purposes
	go func() {
		log.Println(http.ListenAndServe("localhost:"+strconv.Itoa(6060+ID), nil))
	}()

	flag.Parse()

}
