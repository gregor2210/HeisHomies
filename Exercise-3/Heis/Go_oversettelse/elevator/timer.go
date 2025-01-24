package elevator

import (
	"fmt"
	"time"
	"timer"
)

func main() {
	// Start en timer på 2 sekunder
	timer.StartTimer(2.0)

	// Løkke som venter på timeout
	for {
		if timer.HasTimedOut() {
			fmt.Println("Timeren har gått ut!")
			break
		} else {
			fmt.Println("Venter...")
			time.Sleep(500 * time.Millisecond) // Vent litt før du sjekker igjen
		}
	}

	// Stopp timeren
	timer.StopTimer()
	fmt.Println("Timer stoppet.")
}