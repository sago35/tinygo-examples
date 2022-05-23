package main

import (
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// TODO: use atomics to access this value.
var ledColor = []byte{0x00, 0x00, 0x00, 0x00}
var leds = []machine.Pin{machine.LED_RED, machine.LED_GREEN, machine.LED_BLUE, machine.LED}
var hasColorChange = true

var (
	serviceUUID = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
	charUUID    = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
)

func main() {
	println("starting")
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "TinyGo ble-led-server",
	}))
	must("start adv", adv.Start())

	var ledColorCharacteristic bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &ledColorCharacteristic,
				UUID:   charUUID,
				Value:  ledColor[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 4 {
						return
					}
					ledColor[0] = value[0]
					ledColor[1] = value[1]
					ledColor[2] = value[2]
					ledColor[3] = value[3]
					hasColorChange = true
				},
			},
		},
	}))

	for _, led := range leds {
		led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	for {
		for !hasColorChange {
			time.Sleep(10 * time.Millisecond)
		}
		hasColorChange = false
		for i, led := range leds {
			led.Set(ledColor[i] == 0)
		}
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
