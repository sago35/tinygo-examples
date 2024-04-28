//go:build wioterminal

package main

import (
	"machine"
)

func init() {
	button1 = machine.WIO_KEY_C
}
