package connectivity

import "time"

func Connectivity_setup() (chan Worldview_package, <-chan time.Time, chan int) {
	// Staring connectivity threds, and return chanels

	// Setting up TCP connection loop.
	TCP_receive_channel := make(chan Worldview_package)

	// Go routine to send world view every second
	ticker := time.NewTicker(100 * time.Millisecond)
	world_view_send_ticker := ticker.C // Keeps the ticker alive

	// Start TCP receiving in a goroutine
	go TCP_receving_setup(TCP_receive_channel)

	return TCP_receive_channel, world_view_send_ticker, offline_update_chan
}
