//go:build wioterminal
// +build wioterminal

package initialize

import (
	"tinygo.org/x/drivers/examples/ili9341/initdisplay"
	"tinygo.org/x/drivers/ili9341"
)

func Display() *ili9341.Device {
	display := initdisplay.InitDisplay()
	return display
}
