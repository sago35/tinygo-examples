package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

func getCurrentStatus(file string) (bool, time.Time, int, error) {
	//time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	found := false
	var prevTimeOfLap time.Time
	laps := 0

	r, err := os.Open(file)
	if err != nil {
		return found, prevTimeOfLap, laps, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	lastLine := ""
	lastLineWithFound := ""
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Printf("%s\n", lastLine)
		if strings.HasSuffix(lastLine, " found") {
			lastLineWithFound = lastLine
			found = true
		} else {
			found = false
		}
	}

	prevTimeOfLap, err = time.Parse("2006/01/02 15:04:05", lastLineWithFound[:19])
	if err == nil {
		prevTimeOfLap = prevTimeOfLap.Add(-1 * 9 * time.Hour)
	}

	spl := strings.Split(lastLineWithFound, " ")
	x, _ := strconv.ParseUint(spl[3], 10, 64)
	laps = int(x) + 1
	//fmt.Printf("prev : %s\n", prevTimeOfLap.In(time.Local))
	//fmt.Printf("%s\n", lastLineWithFound)

	return found, prevTimeOfLap, laps, nil
}

func run() error {
	found, prevTimeOfLap, laps, err := getCurrentStatus("data.txt")
	if err != nil {
		//return err
	}

	wfh, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer wfh.Close()

	w := io.MultiWriter(os.Stdout, wfh)
	log.SetOutput(w)

	must("enable BLE stack", adapter.Enable())

	// Scan for NUS peripheral.
	timeFromLastFound := time.Now()
	//thresh := 1 * time.Second
	thresh := 3 * time.Minute
	if debug {
		thresh = 1000 * time.Millisecond
	}
	// found -> lost の間は thresh 以上の時間が必要
	// 逆に lost -> found の間も thresh 以上の時間が必要
	for {
		cont := true
		//println("scanning")
		for cont {
			err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
				//fmt.Printf("%#v\n", result.Address.String())
				//fmt.Printf("%#v %d\n", result.Address.String(), -1*time.Until(timer))
				if result.LocalName() != "No 10 - TinyGo LapTimer" {
					if -1*time.Until(timeFromLastFound) > thresh {
						if found {
							log.Printf("lost")
							//time.Sleep(thresh)
						}
						found = false
					}
					return
				}

				// Stop the scan.
				err := adapter.StopScan()
				if err != nil {
					// Unlikely, but we can't recover from this.
					println("failed to stop the scan:", err.Error())
				}
			})
			if err != nil {
				println("could not start a scan:", err.Error())
				return err
			} else {
				if !found {
					lapTime := int(time.Now().Sub(prevTimeOfLap).Seconds())
					log.Printf("%02d:%02d:%02d %d found", lapTime/3600, (lapTime/60)%60, lapTime%60, laps)
					laps++
					prevTimeOfLap = time.Now()
				}
				found = true
				timeFromLastFound = time.Now()
				cont = false
			}
		}
	}
	return nil
}

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", debug, "debug")
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
