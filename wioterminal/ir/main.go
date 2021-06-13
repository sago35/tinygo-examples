package main

// This is the most minimal blinky example and should run almost everywhere.

import (
	"machine"
	"time"
)

func main() {
	led := machine.WIO_IR
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.Low()
		time.Sleep(time.Millisecond * 27)

		led.High()
		time.Sleep(time.Millisecond * 27)
	}
}
