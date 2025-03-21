package connectivity

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

const (
	NumElevators = 3

	// TimeOut for send and receive. If exceeded, connection is lost.
	TimeOut = 3 //seconds

	// Worldview max package size
	MaxPacketSize = 1500

	// UseIPs is set to true if you are gonna use different computer.
	// Remember to set correct ips in IPs
	UseIPs = false
)

var (
	ID int //default 0

	// IPs of all elevators
	IPs = [NumElevators]string{"10.100.23.28", "10.100.23.32", "10.100.23.29"}
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
