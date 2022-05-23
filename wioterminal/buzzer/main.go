package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/tone"
)

func main() {
	bzrPin := machine.WIO_BUZZER
	pwm := machine.TCC0
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	song := []tone.Note{
		tone.C5,
		tone.D5,
		tone.E5,
		tone.F5,
		tone.G5,
		tone.A5,
		tone.B5,
		tone.C6,
		tone.C6,
		tone.B5,
		tone.A5,
		tone.G5,
		tone.F5,
		tone.E5,
		tone.D5,
		tone.C5,
	}

	for {
		for _, val := range song {
			speaker.SetNote(val)
			time.Sleep(time.Second / 2)
		}
	}
}
