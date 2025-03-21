package elevio

func InitIOHandling() (chan ButtonEvent, chan int, chan bool) {
	drvbuttons := make(chan ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	//drv_stop := make(chan bool)

	go PollButtons(drvbuttons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)

	// go PollStopButton(drv_stop) // Ikke implementert enda.

	return drvbuttons, drvFloors, drvObstr

}
