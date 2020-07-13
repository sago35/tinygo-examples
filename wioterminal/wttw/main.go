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

	rtl = initRTL()
	//rtl.Configure(&rtl8720dn.Config{})
	net.ActiveDevice = rtl
	fmt.Printf("connecting\r\n")

	if connectToRTL8720() {
		println("Connected to wifi adaptor.")
		//adaptor.Echo(false)

		connectToAP()
	} else {
		println("")
		failMessage("Unable to connect to wifi adaptor.")
		return
	}

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

	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display.Configure(ili9341.Config{})
	width, height := display.Size()

	display.FillScreen(black)
	backlight.High()

	display.FillRectangle(0, 0, width/2, height/2, white)
	display.FillRectangle(width/2, 0, width/2, height/2, red)
	display.FillRectangle(0, height/2, width/2, height/2, green)
	display.FillRectangle(width/2, height/2, width/2, height/2, blue)
	display.FillRectangle(width/4, height/4, width/2, height/2, black)

	tinyfont.Wrap = true

	loop()
}

func loop() {
	i := 0
	needRedraw := true
	s := tweet.S
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

		if !btnNext.Get() || !btnNext2.Get() {
			i = (i + 1) % len(s)
			needRedraw = true
		} else if !btnPrev.Get() || !btnPrev2.Get() {
			i = (i + len(s) - 1) % len(s)
			needRedraw = true
		} else if !btnPress.Get() {
			fmt.Printf("press\r\n")
			btnPressed()
		}
	}
}

func btnPressed() {
	buf, err := httpGet("http://192.168.1.114:8080/u")
	if err != nil {
		fmt.Printf("err %s\r\n", err.Error())
		return
	}
	t, err := tweet.NewTweet2(buf)
	if err != nil {
		fmt.Printf("err %s\r\n", err.Error())
		return
	}
	fmt.Printf("%#v\r\n", t)

	display.FillScreen(color.RGBA{255, 255, 255, 255})

	tinyfont.WriteLine(display, &TinyFont, 3, 15, t.UserName, color.RGBA{0, 0, 0, 255})
	tinyfont.WriteLine(display, &MplusConst10pt, tinyfont.Cx, tinyfont.Cy, fmt.Sprintf(" @%s\n", t.ScreenName), color.RGBA{158, 158, 158, 255})

	tinyfont.Cy += 3
	tinyfont.WriteLine(display, &MplusConst10pt, 3, tinyfont.Cy, fmt.Sprintf("%s fav:%d rt:%d\n", t.CreatedAt, t.FavoriteCount, t.RetweetCount), color.RGBA{158, 158, 158, 255})
	tinyfont.Cy += 5

	tinyfont.WriteLine(display, &TinyFont, 3, tinyfont.Cy+5, t.FullText, color.RGBA{0, 0, 0, 255})
}
