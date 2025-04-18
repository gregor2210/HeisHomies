package elevio

import "fmt"

func InitIOHandling() (chan ButtonEvent, chan int, chan bool) {
	fmt.Println("Initializing IO handling")
	drvbuttons := make(chan ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	//drv_stop := make(chan bool)

	go PollButtons(drvbuttons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)

	fmt.Println("IO handling initialized")
	// go PollStopButton(drv_stop) // Ikke implementert enda.
	return drvbuttons, drvFloors, drvObstr

}
