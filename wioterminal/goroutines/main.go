package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ili9341"
)

var (
	led1    = machine.LED
	display *ili9341.Device
)

var (
	black = color.RGBA{0, 0, 0, 255}
	gray  = color.RGBA{32, 32, 32, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}
)

func initialize() error {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 48000000,
	})

	backlight := machine.LCD_BACKLIGHT
	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display = ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)
	display.Configure(ili9341.Config{
		Rotation: ili9341.Rotation270,
	})
	display.FillScreen(black)
	backlight.High()

	return nil
}

func errDisp(err error) {
	for {
		fmt.Printf("%s\r\n", err.Error())
		time.Sleep(10 * time.Second)
	}
}

func main() {
	err := initialize()
	if err != nil {
		errDisp(err)
	}

	label1 := NewLabel(240, 0x12)
	label2 := NewLabel(240, 0x12)

	chCnt1 := make(chan uint32, 1)
	chCnt2 := make(chan uint32, 1)

	go timer77ms(chCnt1)
	go timer500ms(chCnt2)

	for {
		select {
		case cnt := <-chCnt1:
			label1.SetText(fmt.Sprintf("timer77ms : %04X", cnt), white)
			display.DrawRGBBitmap(0, 30, label1.buf, label1.w, label1.h)
		case cnt := <-chCnt2:
			label2.SetText(fmt.Sprintf("timer500ms: %04X", cnt), white)
			display.DrawRGBBitmap(0, 50, label2.buf, label2.w, label2.h)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func timer77ms(ch chan<- uint32) {
	cnt := uint32(0)
	for {
		ch <- cnt
		cnt++
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
