package elevio

import "fmt"

func InitIOHandling() (chan ButtonEvent, chan int, chan bool) {
	fmt.Println("Initializing IO handling")
	drvbuttons := make(chan ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)

	go PollButtons(drvbuttons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)

	fmt.Println("IO handling initialized")

	return drvbuttons, drvFloors, drvObstr

}
