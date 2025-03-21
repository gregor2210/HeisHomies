package elevio

func Io_threds_setup() (chan ButtonEvent, chan int, chan bool) {
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	//drv_stop := make(chan bool)

	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)

	// go PollStopButton(drv_stop) // Ikke implementert enda.

	return drv_buttons, drv_floors, drv_obstr

}
