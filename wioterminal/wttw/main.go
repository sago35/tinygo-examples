package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"github.com/sago35/tinygo-examples/wioterminal/wttw/tweet"
	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/rtl8720dn"
	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfont"
)

var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}
)

var (
	sd  sdcard.Device
	rtl *rtl8720dn.Device
)

var (
	led = machine.LED
)

func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display.Configure(ili9341.Config{
		Rotation: ili9341.Rotation180,
	})
	display.FillScreen(white)

	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})
	backlight.High()

	rtl = initRTL()
	//rtl.Configure(&rtl8720dn.Config{})
	net.ActiveDevice = rtl
	fmt.Printf("connecting\r\n")
	tinyfont.WriteLine(display, &TinyFont, 3, 15, "connecting\n", color.RGBA{0, 0, 0, 255})

	if connectToRTL8720() {
		println("Connected to wifi adaptor.")
		tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy, "connected to wifi adaptor.\n", color.RGBA{0, 0, 0, 255})
		//adaptor.Echo(false)

		connectToAP()
	} else {
		println("")
		failMessage("Unable to connect to wifi adaptor.")
		tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy, "unable to connect to wifi adaptor\n", color.RGBA{0, 0, 0, 255})
		return
	}
	tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy, "connected\n", color.RGBA{0, 0, 0, 255})

	tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy, "init sdcard", color.RGBA{0, 0, 0, 255})
	if true {
		sd = sdcard.New(sdSpi, sdCs)
		err := sd.Configure()
		if err != nil {
			fmt.Printf("%s\r\n", err.Error())
			for {
				time.Sleep(time.Hour)
			}
		}
	}
	tinyfont.WriteLine(display, &TinyFont, tinyfont.Cx, tinyfont.Cy, "...done\n", color.RGBA{0, 0, 0, 255})
	time.Sleep(500 * time.Millisecond)

	tinyfont.Wrap = true

	loop()
}

func loop() {
	i := 0
	needRedraw := true
	s := tweet.S

	if true {
		needRedraw = false
		btnPressed(-1)
	}
	for {
		led.Toggle()

		if needRedraw {
			fmt.Printf("redraw\r\n")
			display.FillScreen(color.RGBA{255, 255, 255, 255})
			//tinyfont.WriteLine(display, &TinyFont, 3, 15, "hello あたらしい世界\n憂鬱とかの難しい漢字も問題なしなので普通に twitter から読んできても問題なく表示できる予定。まだ🍺とかは表示できないはず。", color.RGBA{0, 0, 0, 255})
			//tinyfont.WriteLine(display, &TinyFont, 3, 15, "まだ🍺とか😁とかは表示できないはず。", color.RGBA{0, 0, 0, 255})
			//display.DrawFastHLine(0, 239, 15, red)
			tinyfont.WriteLine(display, &TinyFont, 3, 15, fmt.Sprintf("%s", s[i].UserName), color.RGBA{0, 0, 0, 255})
			tinyfont.WriteLine(display, &MplusConst10pt, tinyfont.Cx, tinyfont.Cy, fmt.Sprintf(" @%s\n", s[i].ScreenName), color.RGBA{158, 158, 158, 255})

			tinyfont.Cy += 3
			tinyfont.WriteLine(display, &MplusConst10pt, 3, tinyfont.Cy, fmt.Sprintf("%s fav:%d rt:%d\n", s[i].CreatedAt, s[i].FavoriteCount, s[i].RetweetCount), color.RGBA{158, 158, 158, 255})
			tinyfont.Cy += 5

			tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy+5, s[i].FullText, color.RGBA{0, 0, 0, 255})
			tinyfont.WriteLine(display, &TinyFont, 3, 315, fmt.Sprintf("%d / %d", i+1, len(s)), color.RGBA{64, 64, 64, 255})
			needRedraw = false
		}

		if !btnNext.Get() {
			i = (i + 1) % len(s)
			needRedraw = true
		} else if !btnNext2.Get() {
			// down
			btnPressed(-1)
		} else if !btnPrev.Get() {
			i = (i + len(s) - 1) % len(s)
			needRedraw = true
		} else if !btnPrev2.Get() {
			// up
			btnPressed(1)
		} else if !btnPress.Get() {
			fmt.Printf("press\r\n")
			btnPressed(1)
		}
	}
}

var bmp [240 * 10]uint16
var prevID int64

func btnPressed(offset int) {
	url := "http://192.168.1.114:8081/u/"
	if prevID == 0 {
	} else if offset < 0 {
		url += fmt.Sprintf("?max=%d", prevID)
	} else {
		url += fmt.Sprintf("?since=%d", prevID)
	}
	buf, err := httpGet(url)
	if err != nil {
		fmt.Printf("err %s\r\n", err.Error())
		// reconnect
		rtl.Configure(&rtl8720dn.Config{})
		if connectToRTL8720() {
			println("Connected to wifi adaptor.")
			//adaptor.Echo(false)

			connectToAP()
		} else {
			println("")
			failMessage("Unable to connect to wifi adaptor.")
			return
		}
		return
	}
	t, err := tweet.NewTweet2(buf)
	if err != nil {
		fmt.Printf("err %s\r\n", err.Error())
		return
	}
	if prevID == t.Id {
		// skip
		m.Msgbox("新しい tweet はありません", 0, 120)
		time.Sleep(1 * time.Second)
	}
	display.FillScreen(color.RGBA{255, 255, 255, 255})

	prevID = t.Id
	fmt.Printf("%#v\r\n", t)

	tinyfont.WriteLine(display, &TinyFont, 3, 15, t.UserName, color.RGBA{0, 0, 0, 255})
	tinyfont.WriteLine(display, &MplusConst10pt, tinyfont.Cx, tinyfont.Cy, fmt.Sprintf(" @%s\n", t.ScreenName), color.RGBA{158, 158, 158, 255})

	tinyfont.Cy += 3
	tinyfont.WriteLine(display, &MplusConst10pt, 3, tinyfont.Cy, fmt.Sprintf("%s fav:%d rt:%d\n", t.CreatedAt, t.FavoriteCount, t.RetweetCount), color.RGBA{158, 158, 158, 255})
	tinyfont.Cy += 5

	tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy+5, t.FullText, color.RGBA{0, 0, 0, 255})
}
