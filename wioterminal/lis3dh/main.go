// Connects to a LIS3DH I2C accelerometer on the Adafruit Circuit Playground Express.
package main

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/lis3dh"
)

var i2c = machine.I2C0

func main() {
	i2c.Configure(machine.I2CConfig{SCL: machine.SCL0_PIN, SDA: machine.SDA0_PIN})

	accel := lis3dh.New(i2c)
	accel.Address = lis3dh.Address0 // address on the Wio Terminal
	accel.Configure()
	accel.SetRange(lis3dh.RANGE_2_G)

	println(accel.Connected())

	for {
		x, y, z, _ := accel.ReadAcceleration()
		fmt.Printf("X: %-10d Y: %-10d Z: %-10d\r\n", x, y, z)

		//rx, ry, rz := accel.ReadRawAcceleration()
		//println("X (raw):", rx, "Y (raw):", ry, "Z (raw):", rz)

		time.Sleep(time.Millisecond * 100)
	}
}
