package main

import (
	"device/arm"
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/notosans"
)

var timerCh = make(chan struct{}, 1)

var (
	led = machine.LED
	d0  = machine.D0
	d11 = machine.D11
	d12 = machine.D12
	a2  = machine.A2
	x   = int16(0)
	y   = int16(0)
)

var (
	currentTime uint64
	prevTime    uint64
	prevTime2   uint64
	laps        uint32
)

func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	d0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	d11.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d12.Configure(machine.PinConfig{Mode: machine.PinOutput})

	a2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})

	// timer fires 500 times per second
	arm.SetupSystemTimer(machine.CPUFrequency() / 500)

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32,
		Width:   128,
		Height:  32,
	})

	display.ClearDisplay()

	clear := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	ticker := uint32(0)

	if false {
		// for debug
		white = color.RGBA{0, 0, 0, 255}
		clear = color.RGBA{255, 255, 255, 255}
	}

	for {
		<-timerCh

		if (ticker & 0x3F) == 0 {
			display.SetPixel(x, 0, clear)
			display.SetPixel(x, 2, clear)
			display.SetPixel(x, y, white)

			//tinyfont.WriteLine(&display, &notosans.Notosans12pt, 3, 15, "hello", clear)
			for yy := int16(5 + 15*0); yy < 16+15*0; yy++ {
				for xx := int16(0); xx < 60; xx++ {
					display.SetPixel(xx, yy, clear)
				}
			}
			tinyfont.WriteLine(&display, &notosans.Notosans12pt, 3, 15, fmt.Sprintf("%02d:%02d:%02d", currentTime/1000/60, (currentTime/1000)%60, (currentTime%1000)/10), white)

			for yy := int16(5 + 15*1); yy < 16+15*1; yy++ {
				for xx := int16(0); xx < 60; xx++ {
					display.SetPixel(xx, yy, clear)
				}
			}

			tinyfont.WriteLine(&display, &notosans.Notosans12pt, 3, 30, fmt.Sprintf("%d laps", laps), white)

			if prevTime != 0 {
				zz := int16(0)
				for zz = 0; zz < 2; zz++ {
					for yy := int16(5 + 15*zz); yy < 17+15*zz; yy++ {
						for xx := int16(64); xx < 64+60; xx++ {
							display.SetPixel(xx, yy, clear)
						}
					}
				}
				tinyfont.WriteLine(&display, &notosans.Notosans12pt, 3+64, 15, fmt.Sprintf("%02d:%02d:%02d", prevTime/1000/60, (prevTime/1000)%60, (prevTime%1000)/10), white)
				if prevTime2 != 0 {
					tinyfont.WriteLine(&display, &notosans.Notosans12pt, 3+64, 30, fmt.Sprintf("%02d:%02d:%02d", prevTime2/1000/60, (prevTime2/1000)%60, (prevTime2%1000)/10), white)
				}
				prevTime2 = prevTime
				prevTime = 0
				laps++
			}
			display.Display()

			//led.Toggle()
			d12.Toggle()
			x = (x + 1) % 128
		}
		ticker++
	}
}

//export SysTick_Handler
func timer_isr() {
	tick()
	tick2()

	select {
	case timerCh <- struct{}{}:
	default:
		// The consumer is running behind.
	}
}

const (
	stIDL = iota
	stIDL2EXPIRE
	stEXPIRE
	stEXPIRED
)

var (
	state = stIDL
	cnt   = 0
	cntB  = 0
	prev  = false
	guard = 0
)

func tick() {
	current := a2.Get()
	currentTime += 2

	led.Toggle()
	if current {
		y = 2
		//led.High()
	} else {
		y = 0
		//led.Low()
	}

	if prev != current {
		cnt++
		cntB = 0
	} else {
		cntB++
	}
	prev = current

	switch state {
	case stIDL:
		if 10 < cnt {
			state = stEXPIRE
		}
		if 500 < cntB {
			cnt = 0
		}
	case stEXPIRE:
		//d11.High()
		guard = 10 * 500
		state = stEXPIRED
		prevTime = currentTime
		currentTime = 0
	case stEXPIRED:
		guard--
		if 0 == guard {
			cnt = 0
			//d11.Low()
			state = stIDL
		}
	}
}

const (
	tick2stWaitStartBitLow = iota
	tick2stWaitStartBitHigh
	tick2stParse
	tick2stWaitLow
	tick2stWaitHigh
)

var (
	tick2currentTime uint32
	tick2prevTime    uint32
	tick2state       = tick2stWaitStartBitLow
	tick2cnt         uint32
	buf              [10]byte
)

// tick2 is called per 2ms
func tick2() {
	current := d0.Get()
	tick2currentTime += 2
	idx := 0
	bitIdx := byte(0)

	switch tick2state {
	case tick2stWaitStartBitLow:
		if !current {
			d11.Toggle()
			tick2cnt += 2
			if 5 < tick2cnt {
				tick2state = tick2stWaitStartBitHigh
				tick2cnt = 0
			}
		} else {
			tick2cnt = 0
		}
	case tick2stWaitStartBitHigh:
		if current {
			d11.Low()
			tick2state = tick2stParse
			tick2cnt = 0
		}
	case tick2stParse:
		// blocking 処理 (数 ms 以上)
		t := 0
		for d0.Get() {
			// 初回の Low の手前までを読み飛ばす
		}

		for {
			d11.Toggle()
			for !d0.Get() {
				// High に戻るのを待つ
			}
			t = time.Now().Nanosecond()
			for d0.Get() {
				// 信号の終わり (Low になる) を待つ
				if 12*1000*1000 < time.Now().Nanosecond()-t {
					// 12ms 以上 High が継続する場合は抜ける
					if bitIdx == 0 && 0 < idx {
						idx--
					}
					//fmt.Printf("-= %v\r\n", buf[:idx])
					fmt.Print(fmt.Sprintf("-= %02X %02X %02X %02X %02X \r\n", buf[0], buf[1], buf[2], buf[3], buf[4]))
					//fmt.Printf("-= %02X %02X %02X %02X %02X \r\n", buf[0], buf[1], buf[2], buf[3], buf[4])
					for i := range buf {
						buf[i] = 0x00
					}
					tick2state = tick2stWaitStartBitLow
					return
				}
			}
			t = time.Now().Nanosecond() - t

			if 10*1000*1000 <= t {
				// stop data
				break
			} else if 1*1000*1000 <= t {
				// high
				buf[idx] |= 1 << (7 - bitIdx)
			} else {
				// low
			}

			bitIdx++
			if 8 <= bitIdx {
				idx++
				bitIdx = 0
			}
		}
		d11.Toggle()
		if bitIdx == 0 && 0 < idx {
			idx--
		}
		fmt.Printf("-- %v\r\n", buf[:1])

		tick2state = tick2stWaitStartBitLow
	case tick2stWaitLow:
	case tick2stWaitHigh:
	}
}
