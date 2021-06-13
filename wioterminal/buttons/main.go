package main

import (
	"fmt"
	"machine"
	"time"
)

const (
	led = machine.LED
)

func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.WIO_KEY_A.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_KEY_B.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_KEY_C.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	machine.WIO_5S_UP.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_LEFT.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_RIGHT.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_DOWN.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_PRESS.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	for {
		if !machine.WIO_KEY_A.Get() {
			led.Low()
			fmt.Printf("machine.WIO_KEY_A pressed\r\n")
		} else if !machine.WIO_KEY_B.Get() {
			led.Low()
			fmt.Printf("machine.WIO_KEY_B pressed\r\n")
		} else if !machine.WIO_KEY_C.Get() {
			led.Low()
			fmt.Printf("machine.WIO_KEY_C pressed\r\n")
		} else if !machine.WIO_5S_UP.Get() {
			led.Low()
			fmt.Printf("machine.WIO_5S_UP pressed\r\n")
		} else if !machine.WIO_5S_LEFT.Get() {
			led.Low()
			fmt.Printf("machine.WIO_5S_LEFT pressed\r\n")
		} else if !machine.WIO_5S_RIGHT.Get() {
			led.Low()
			fmt.Printf("machine.WIO_5S_RIGHT pressed\r\n")
		} else if !machine.WIO_5S_DOWN.Get() {
			led.Low()
			fmt.Printf("machine.WIO_5S_DOWN pressed\r\n")
		} else if !machine.WIO_5S_PRESS.Get() {
			led.Low()
			fmt.Printf("machine.WIO_5S_PRESS pressed\r\n")
		} else {
			led.High()
		}

		time.Sleep(time.Millisecond * 10)
	}
}
