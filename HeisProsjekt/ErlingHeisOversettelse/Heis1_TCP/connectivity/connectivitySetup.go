package connectivity

import (
	"fmt"
	"time"
)

func ConnectivitySetup() (chan WorldviewPackage, <-chan time.Time, chan int) {
	fmt.Println("Setting up connectivity")
	tcpReceiveChannel := make(chan WorldviewPackage)
	go TcpReceivingSetup(tcpReceiveChannel)

	// Go routine to send world view every second
	ticker := time.NewTicker(100 * time.Millisecond)
	worldViewSendTicker := ticker.C // Keeps the ticker alive

	fmt.Println("Connectivity set up")
	return tcpReceiveChannel, worldViewSendTicker, offlineUpdateChan
}
