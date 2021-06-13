package main

import (
	"machine"

	"tinygo.org/x/drivers/flash"
)

func main() {
	RunFor(
		flash.NewQSPI(
			machine.QSPI_CS,
			machine.QSPI_SCK,
			machine.QSPI_DATA0,
			machine.QSPI_DATA1,
			machine.QSPI_DATA2,
			machine.QSPI_DATA3,
		),
	)
}
