package connectivity

import "time"

func ConnectivitySetup() (chan WorldviewPackage, <-chan time.Time, chan int) {

	tcpReceiveChannel := make(chan WorldviewPackage)
	go TcpReceivingSetup(tcpReceiveChannel)

	// Go routine to send world view every second
	ticker := time.NewTicker(100 * time.Millisecond)
	worldViewSendTicker := ticker.C // Keeps the ticker alive

	return tcpReceiveChannel, worldViewSendTicker, offlineUpdateChan
}
