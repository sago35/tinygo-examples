package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	d0 := machine.D0
	d0.Configure(machine.PinConfig{Mode: machine.PinOutput})

	d8 := machine.D8
	d8.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	d9 := machine.D9
	d9.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	uart := machine.UART1
	uart.Configure(machine.UARTConfig{
		BaudRate: 34000,
		TX:       machine.UART_TX_PIN,
		RX:       machine.UART_RX_PIN,
	})

	go func() {
		for {
			if !d8.Get() {
				volumeUp(uart)
				time.Sleep(200 * time.Millisecond)
			}
			if !d9.Get() {
				volumeDown(uart)
				time.Sleep(200 * time.Millisecond)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	for {
		led.High()
		time.Sleep(100 * time.Millisecond)
		led.Low()
		time.Sleep(100 * time.Millisecond)
	}
}

var (
	buf [4096]byte
)

func irEncode(in []byte) []byte {
	copy(buf[8*0:], []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	copy(buf[8*1:], []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	copy(buf[8*2:], []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	copy(buf[8*3:], []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	copy(buf[8*4:], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	copy(buf[8*5:], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	start := 8 * 6

	for _, b := range in {
		for i := 0; i < 8; i++ {
			bit := b & 0x80

			if bit == 0 {
				copy(buf[start:], []byte{0x00, 0x00, 0xFF, 0xFF})
				start += 4
			} else {
				copy(buf[start:], []byte{0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
				start += 8
			}

			b <<= 1
		}
	}

	return buf[:start]
}

func volumeUp(uart machine.UART) {
	uart.Write(irEncode([]byte{0x20, 0xDF, 0x40, 0xBF, 0x00}))
}

func volumeDown(uart machine.UART) {
	uart.Write(irEncode([]byte{0x20, 0xDF, 0xC0, 0x3F, 0x00}))
}
