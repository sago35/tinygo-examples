package main

import (
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var ledColor = []byte{0x00, 0x00, 0x00, 0x00}

var (
	serviceUUID = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
	charUUID    = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
)

func main() {
	println("starting")
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "No 10 - TinyGo LapTimer",
		Interval:  bluetooth.NewDuration(32 * time.Millisecond),
	}))
	must("start adv", adv.Start())

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	//var err error
	for {
		//fmt.Printf("start\r\n")
		//err = adv.Start()
		//if err != nil {
		//	fmt.Printf("err %s\r\n", err.Error())
		//}

		led.Low()
		time.Sleep(time.Millisecond * 100)
		led.High()
		time.Sleep(time.Millisecond * 100)
		//fmt.Printf("stop\r\n")
		//err = adv.Stop()
		//if err != nil {
		//	fmt.Printf("err %s\r\n", err.Error())
		//}
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
