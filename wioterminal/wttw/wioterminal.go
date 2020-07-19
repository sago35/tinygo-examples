// +build wioterminal

package main

import (
	"machine"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/rtl8720dn"
)

var (
	display = ili9341.NewSpi(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)

	backlight = machine.LCD_BACKLIGHT
)

func init() {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		MOSI:      machine.LCD_MOSI_PIN,
		MISO:      machine.LCD_MISO_PIN,
		Frequency: 48000000,
	})
}

// sdcard
var (
	sdSpi machine.SPI
	sdCs  machine.Pin
)

func init() {
	sdSpi = machine.SPI2
	sdSpi.Configure(machine.SPIConfig{
		SCK:       machine.SCK2,
		MOSI:      machine.MOSI2,
		MISO:      machine.MISO2,
		Frequency: 24000000,
		LSBFirst:  false,
		Mode:      0, // phase=0, polarity=0
	})

	sdCs = machine.SS2
}

// buttons
var (
	btnNext = machine.BUTTON_3
	btnPrev = machine.BUTTON_2

	btnNext2 = machine.WIO_5S_DOWN
	btnPrev2 = machine.WIO_5S_UP

	btnPress = machine.WIO_5S_PRESS
)

func init() {
	btnNext.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnPrev.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnNext2.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnPrev2.Configure(machine.PinConfig{Mode: machine.PinInput})

	btnPress.Configure(machine.PinConfig{Mode: machine.PinInput})
}

// rtl8720 wifi
func initRTL() *rtl8720dn.Device {
	spi := machine.SPI1
	spi.Configure(machine.SPIConfig{
		SCK:       machine.SCK1,
		MOSI:      machine.MOSI1,
		MISO:      machine.MISO1,
		Frequency: 40000000,
		LSBFirst:  false,
		Mode:      0, // phase=0, polarity=0
	})

	chipPu := machine.RTL8720D_CHIP_PU
	syncPin := machine.RTL8720D_GPIO0
	csPin := machine.SS1
	uartRxPin := machine.UART2_RX_PIN

	rtl := rtl8720dn.New(spi, chipPu, syncPin, csPin, uartRxPin)
	rtl.Configure(&rtl8720dn.Config{})
	return rtl
}
