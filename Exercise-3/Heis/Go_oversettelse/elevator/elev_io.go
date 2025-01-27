package elevator

type ElevInputDevice struct {
	FloorSensor     func() int
	RequestButton   func(int, Button) int
	StopButton      func() int
	Obstruction     func() int
}

func GetInputDevice() ElevInputDevice {
	return ElevInputDevice{
		FloorSensor:     hardwareGetFloorSensorSignal,
		RequestButton:   wrapRequestButton,
		StopButton:      hardwareGetStopSignal,
		Obstruction:     hardwareGetObstructionSignal,
	}
}

type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, Button, int)
	DoorLight          func(int)
	StopButtonLight    func(int)
	MotorDirection     func(Dirn)
}

func GetOutputDevice() ElevOutputDevice {
	return ElevOutputDevice{
		FloorIndicator:     hardwareSetFloorIndicator,
		RequestButtonLight: wrapRequestButtonLight,
		DoorLight:          hardwareSetDoorOpenLamp,
		StopButtonLight:    hardwareSetStopLamp,
		MotorDirection:     wrapMotorDirection,
	}
}