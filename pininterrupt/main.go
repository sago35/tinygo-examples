package main

import (
	"fmt"
	"machine"
	"runtime/volatile"
	"time"

	"github.com/sago35/tinygo-examples/pininterrupt/st7032"
)

/*
PyPortal
D3  : Switch (with pull-up)
D4  : LED
I2C : AE-AQM0802A (I2C 8 x 2 character LCD with st7032)
*/

var (
	led1    = machine.LED
	led2    = machine.D3
	button  = machine.D4
	aqm0802 = st7032.New(&machine.I2C0, 0x3E)
)

// scrollp is a function to scroll str
func scrollp(d *st7032.Device, lineno uint8, str string) {
	current := `        `
	for i := 0; i < len(str); i++ {
		aqm0802.SetCursor(0, lineno)
		aqm0802.Print(current)
		time.Sleep(200 * time.Millisecond)
		current = current[1:] + string(str[i])
	}
	for i := 0; i < 9; i++ {
		aqm0802.SetCursor(0, lineno)
		aqm0802.Print(current)
		time.Sleep(250 * time.Millisecond)
		current = current[1:] + ` `
	}
}

// lv is used to create the first line of the LCD.
type lv byte

func (d lv) String() string {
	buf := make([]byte, 8)
	for i := 0; i < 8; i++ {
		if (d & 0x80) > 0 {
			buf[i] = byte('\xFF')
		} else {
			buf[i] = byte('_')
		}
		d <<= 1
	}
	return string(buf)
}

func (d *lv) Add(b bool) {
	if b {
		*d <<= 1
		*d += 1
	} else {
		*d <<= 1
	}
}

// pinInterruptChan notifies an interrupt to a pin via chan bool.
// p.Get() can get the button state, but the button state at the time of getting it may be old.
// Therefore, variables of type volatile.Register8 are used to manage them in the function.
func pinInterruptChan(pin machine.Pin, state *volatile.Register8) <-chan bool {
	ch := make(chan bool, 3)

	pin.SetInterrupt(machine.PinToggle, func(p machine.Pin) {
		b := false
		if state.Get() != 1 {
			state.Set(1)
			b = true
		} else {
			state.Set(0)
		}
		select {
		case ch <- b:
		default:
		}
	})

	return ch
}

func main() {
	chCnt1 := make(chan uint32)
	chCnt2 := make(chan uint32, 1)
	chDisp := make(chan lv, 3)

	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})

	button.Configure(machine.PinConfig{Mode: machine.PinInput})
	var btnState volatile.Register8
	chBtn := pinInterruptChan(button, &btnState)

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})

	aqm0802.Configure()
	aqm0802.SetContrast(5)

	scrollp(aqm0802, 0, "TinyGo demo")
	time.Sleep(500 * time.Millisecond)
	aqm0802.Clear()

	go timer77ms(chCnt1, led2)
	go timer500ms(chCnt2)
	go disp(chBtn, chDisp, led1)

	for {
		select {
		case d := <-chDisp:
			aqm0802.SetCursor(0, 0)
			aqm0802.Print(d.String())

		case cnt := <-chCnt1:
			aqm0802.SetCursor(0, 1)
			aqm0802.Print(fmt.Sprintf("%4d", cnt))

		case cnt := <-chCnt2:
			aqm0802.SetCursor(4, 1)
			aqm0802.Print(fmt.Sprintf("%4d", cnt))
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func timer77ms(ch chan<- uint32, led machine.Pin) {
	cnt := uint32(0)
	for {
		ch <- cnt
		cnt++

		led.Toggle()
		time.Sleep(77 * time.Millisecond)
	}
}

func timer500ms(ch chan<- uint32) {
	cnt := uint32(0)
	for {
		ch <- cnt
		cnt++

		time.Sleep(500 * time.Millisecond)
	}
}

func disp(ch <-chan bool, chDisp chan<- lv, led machine.Pin) {
	s := time.Now()
	e := time.Now()
	d := lv(0xFF)
	btn := false
	for {
		select {
		case btn = <-ch:
			led.Set(btn)
		default:
		}

		e = time.Now()
		if e.UnixNano()-s.UnixNano() > 100*1000*1000 {
			d.Add(btn)
			chDisp <- d
			s = e
		}
		time.Sleep(1 * time.Millisecond)
	}
}
