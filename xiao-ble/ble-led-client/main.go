package main

// This is the most minimal blinky example and should run almost everywhere.

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

// bluetooth
var (
	adapter     = bluetooth.DefaultAdapter
	dc          bluetooth.DeviceCharacteristic
	serviceUUID = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
	charUUID    = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
)

var (
	patterns = [][]byte{
		{0xFF, 0x00, 0x00, 0x00},
		{0x00, 0xFF, 0x00, 0x00},
		{0x00, 0x00, 0xFF, 0x00},
	}
)

func initBLE() {
	println("starting")
	must("enable BLE stack", adapter.Enable())

	// The address to connect to. Set during scanning and read afterwards.
	var foundDevice bluetooth.ScanResult

	// Scan for NUS peripheral.
	println("Scanning...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		fmt.Printf("%#v\n", result.Address.String())
		if result.LocalName() != "TinyGo ble-led-server" {
			return
		}
		foundDevice = result

		// Stop the scan.
		err := adapter.StopScan()
		if err != nil {
			// Unlikely, but we can't recover from this.
			println("failed to stop the scan:", err.Error())
		}
	})
	if err != nil {
		println("could not start a scan:", err.Error())
		return
	}

	// Found a device: print this event.
	if name := foundDevice.LocalName(); name == "" {
		print("Connecting to ", foundDevice.Address.String(), "...")
		println()
	} else {
		print("Connecting to ", name, " (", foundDevice.Address.String(), ")...")
		println()
	}

	// Found a NUS peripheral. Connect to it.
	device, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		println("Failed to connect:", err.Error())
		return
	}

	// Connected. Look up the Nordic UART Service.
	println("Discovering service...")
	services, err := device.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil {
		println("Failed to discover the Nordic UART Service:", err.Error())
		return
	}
	service := services[0]

	// Get the two characteristics present in this service.
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{charUUID})
	if err != nil {
		println("Failed to discover RX and TX characteristics:", err.Error())
		return
	}

	dc = chars[0]
}

func main() {
	initBLE()
	ledValue := []byte{0x00, 0x00, 0x00, 0x00}
	index := 0
	for {
		ledValue[0] = patterns[index][0]
		ledValue[1] = patterns[index][1]
		ledValue[2] = patterns[index][2]
		index = (index + 1) % len(patterns)
		setLEDs(ledValue)
		time.Sleep(200 * time.Millisecond)
	}
}

func setLEDs(b []byte) {
	dc.WriteWithoutResponse(b)
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
