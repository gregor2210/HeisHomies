package connectivity

import "time"

func Connectivity_setup() (chan Worldview_package, <-chan time.Time, chan int) {

	TCP_receive_channel := make(chan Worldview_package)
	go TCP_receving_setup(TCP_receive_channel)

	// Go routine to send world view every second
	ticker := time.NewTicker(100 * time.Millisecond)
	world_view_send_ticker := ticker.C // Keeps the ticker alive

	return TCP_receive_channel, world_view_send_ticker, offline_update_chan
}
