package elevio

func Io_threds_setup() (chan ButtonEvent, chan int, chan bool) {
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	//drv_stop := make(chan bool)

	go PollButtons(drv_buttons)
	//søker igjennom alle etasjene og sjekker alle typer knapper for den etasjen.
	//Den sjekker ved å sende en tpc getbutton(etasje, knappetype) og får tilbake true/false. Dersom dette er anderledes enn fra forigje gang den sjekket og den nå er nå true.
	//Skriver den til chanelen drv_buttons. Den skriver da en ButtonEvent (struct) med etasje og knappetype.
	go PollFloorSensor(drv_floors)
	//Sjekker om heisen er i en etasje og den etasjen er ulik det den var sist gang den sjekket. Dersom den er det, skriver den til chanelen drv_floors. Den skriver da etasjen heisen er i, i form av en int.
	go PollObstructionSwitch(drv_obstr)

	//Sjekker om det er en obstruction i heisen. Dersom statuesn på obstruction endrer seg så skriver den true til chanelen drv_obstr. Dersom det ikke er obstruction, skriver den false.
	// go PollStopButton(drv_stop) // Ikke implementert enda.

	return drv_buttons, drv_floors, drv_obstr

}
