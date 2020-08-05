package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/buzzer"
	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/lis3dh"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}
)

var ()

func init() {
}

func main() {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 48000000,
	})

	machine.InitADC()

	backlight := machine.LCD_BACKLIGHT
	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display := ili9341.NewSpi(
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

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ir := machine.WIO_IR
	ir.Configure(machine.PinConfig{Mode: machine.PinOutput})

	btnA := machine.WIO_KEY_A
	btnA.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnB := machine.WIO_KEY_B
	btnB.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnC := machine.WIO_KEY_C
	btnC.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnU := machine.WIO_5S_UP
	btnU.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnL := machine.WIO_5S_LEFT
	btnL.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnR := machine.WIO_5S_RIGHT
	btnR.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnD := machine.WIO_5S_DOWN
	btnD.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnP := machine.WIO_5S_PRESS
	btnP.Configure(machine.PinConfig{Mode: machine.PinInput})

	lightSensor := machine.ADC{Pin: machine.WIO_LIGHT}
	lightSensor.Configure()

	mic := machine.ADC{Pin: machine.WIO_MIC}
	mic.Configure()

	label1 := NewLabel(240, 0x12)
	label2 := NewLabel(240, 0x12)

	machine.I2C1.Configure(machine.I2CConfig{SCL: machine.SCL1_PIN, SDA: machine.SDA1_PIN})
	accel := lis3dh.New(machine.I2C1)
	accel.Address = lis3dh.Address0
	accel.Configure()
	accel.SetRange(lis3dh.RANGE_2_G)

	label3 := NewLabel(320, 0x12)
	label4 := NewLabel(160, 0x12)
	label5 := NewLabel(160, 0x12)

	bzrPin := machine.WIO_BUZZER
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bzr := buzzer.New(bzrPin)

	cnt := 0
	lastRequestTime := time.Now()
	rawXY := [16][2]int16{}
	rawXYi := 0

	for {
		if btnC.Get() {
			display.FillRectangle(10, 10, 10, 10, white)
		} else {
			display.FillRectangle(10, 10, 10, 10, red)
		}

		if btnB.Get() {
			display.FillRectangle(30, 10, 10, 10, white)
		} else {
			display.FillRectangle(30, 10, 10, 10, red)
		}

		if btnA.Get() {
			display.FillRectangle(50, 10, 10, 10, white)
		} else {
			display.FillRectangle(50, 10, 10, 10, red)
		}

		if btnU.Get() {
			display.FillRectangle(280, 180, 10, 10, white)
		} else {
			display.FillRectangle(280, 180, 10, 10, red)
		}

		if btnP.Get() {
			display.FillRectangle(280, 200, 10, 10, white)
		} else {
			display.FillRectangle(280, 200, 10, 10, red)
			bzr.Tone(buzzer.G7, 0.02)
		}

		if btnD.Get() {
			display.FillRectangle(280, 220, 10, 10, white)
		} else {
			display.FillRectangle(280, 220, 10, 10, red)
		}

		if btnL.Get() {
			display.FillRectangle(260, 200, 10, 10, white)
		} else {
			display.FillRectangle(260, 200, 10, 10, red)
		}

		if btnR.Get() {
			display.FillRectangle(300, 200, 10, 10, white)
		} else {
			display.FillRectangle(300, 200, 10, 10, red)
		}

		label1.SetText(fmt.Sprintf("WIO_LIGHT : %04X", lightSensor.Get()), white)
		display.DrawRGBBitmap(0, 30, label1.buf, label1.w, label1.h)

		label2.SetText(fmt.Sprintf("WIO_MIC   : %04X", mic.Get()), white)
		display.DrawRGBBitmap(0, 50, label2.buf, label2.w, label2.h)

		x, y, z, _ := accel.ReadAcceleration()
		label3.SetText(fmt.Sprintf("LIS3DH    : X=%8d", x), white)
		display.DrawRGBBitmap(0, 80, label3.buf, label3.w, label3.h)
		label4.SetText(fmt.Sprintf(": Y=%8d", y), white)
		display.DrawRGBBitmap(0x0b*10, 100, label4.buf, label4.w, label4.h)
		label5.SetText(fmt.Sprintf(": Z=%8d", z), white)
		display.DrawRGBBitmap(0x0b*10, 120, label5.buf, label5.w, label5.h)

		x2, y2, _ := accel.ReadRawAcceleration()
		for i := range rawXY {
			display.FillRectangle(50-1-rawXY[(rawXYi+i)%16][1], 170-1+rawXY[(rawXYi+i)%16][0], 2, 2, black)
		}
		rawXY[rawXYi][0] = x2 / 256
		rawXY[rawXYi][1] = y2 / 256

		display.DrawFastVLine(50, 120, 220, white)
		display.DrawFastHLine(0, 100, 170, white)

		for i := range rawXY {
			i = 15 - i
			display.FillRectangle(50-1-rawXY[(rawXYi+i)%16][1], 170-1+rawXY[(rawXYi+i)%16][0], 2, 2, color.RGBA{255 - 16*(uint8((i)%16)), 0, 0, 255})
		}
		rawXYi = (rawXYi + 16 - 1) % 16

		for time.Now().Sub(lastRequestTime).Milliseconds() <= 50 {
		}
		lastRequestTime = time.Now()

		cnt++
		if (cnt & 0x0F) == 0 {
			led.Toggle()
			ir.Toggle()
		}
	}
}

type label struct {
	buf        []uint16
	w          int16
	h          int16
	fontHeight int16
}

func NewLabel(w, h int) *label {
	return &label{
		buf:        make([]uint16, w*h),
		w:          int16(w),
		h:          int16(h),
		fontHeight: int16(tinyfont.GetGlyph(&freemono.Regular9pt7b, '0').Height),
	}
}

func (l *label) Size() (int16, int16) {
	return l.w, l.h
}

func (l *label) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || y < 0 || l.w < x || l.h < y {
		return
	}
	l.buf[y*l.w+x] = ili9341.RGBATo565(c)
}

func (l *label) Display() error {
	return nil
}

func (l *label) SetText(str string, c color.RGBA) {
	for i := range l.buf {
		l.buf[i] = 0
	}

	tinyfont.WriteLine(l, &freemono.Regular9pt7b, 3, l.fontHeight, str, c)
}
