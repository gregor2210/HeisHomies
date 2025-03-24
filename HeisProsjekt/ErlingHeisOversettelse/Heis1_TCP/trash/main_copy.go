package main

import (
	"Driver-go/elevio"
	"fmt"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MotorUp
	//elevio.SetMotorDirection(d)

	drvbuttons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drvbuttons)
	//søker igjennom alle etasjene og sjekker alle typer knapper for den etasjen.
	//Den sjekker ved å sende en tpc getbutton(etasje, knappetype) og får tilbake true/false. Dersom dette er anderledes enn fra forigje gang den sjekket og den nå er nå true.
	//Skriver den til chanelen drvbuttons. Den skriver da en ButtonEvent (struct) med etasje og knappetype.
	go elevio.PollFloorSensor(drvFloors)
	//Sjekker om heisen er i en etasje og den etasjen er ulik det den var sist gang den sjekket. Dersom den er det, skriver den til chanelen drvFloors. Den skriver da etasjen heisen er i, i form av en int.
	go elevio.PollObstructionSwitch(drvObstr)
	//Sjekker om det er en obstruction i heisen. Dersom statuesn på obstruction endrer seg så skriver den true til chanelen drvObstr. Dersom det ikke er obstruction, skriver den false.
	go elevio.PollStopButton(drv_stop)
	//Sjekker om stopknappen er trykket inn. Den vil skrive true eller false til chanelen drv_stop når statusen endrer seg.

	fmt.Println("Started!")

	//inputPollRateMs := 25
	//elevator := fsm.Elevator{ //lager en ny heis
	//e := fsm.NewElevator()

	//}

	for {
		select {
		case a := <-drvbuttons:
			fmt.Printf("%+v\n", a)
			//elevio.SetButtonLamp(a.Button, a.Floor, true)
			//fsm.Fsm_button_clicked_selecter(a.Button, a.Floor)

		case a := <-drvFloors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevio.MotorDown
			} else if a == 0 {
				d = elevio.MotorUp
			}
			elevio.SetMotorDirection(d)

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MotorStop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}

}
