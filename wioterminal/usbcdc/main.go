package main

import (
	"fmt"
	"strconv"
	"strings"

	"machine"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	select {}
}

func run() error {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	lcd := machine.LCD_BACKLIGHT
	lcd.Configure(machine.PinConfig{Mode: machine.PinOutput})
	lcd.Low()

	ser := machine.Serial

	server := CommandServer{
		LED: led,
		LCD: lcd,
		Messages: []string{
			"hello world",
			"こんにちは 世界",
		},
	}

	i := 0
	input := [256]byte{}
	for {
		if ser.Buffered() > 0 {
			data, _ := ser.ReadByte()

			switch data {
			case 10, 13: // CR or LF
				if i == 0 {
					continue
				}

				ret, err := server.RunCommand(string(input[:i]))
				if err != nil {
					return err
				}
				fmt.Print(ret)
				i = 0
			default:
				input[i] = data
				i++
			}
		}
	}

	return nil
}

type CommandServer struct {
	LED      machine.Pin
	LCD      machine.Pin
	Messages []string
}

func (s *CommandServer) RunCommand(command string) (string, error) {
	if command == "" {
		return "ok\r\n", nil
	}

	spl := strings.Split(command, " ")
	switch spl[0] {
	case "led", "lcd":
		target := s.LED
		if spl[0] == "lcd" {
			target = s.LCD
		}

		if len(spl) == 1 {
			target.Toggle()
		} else if spl[1] == "on" {
			target.High()
		} else {
			target.Low()
		}
	case "msg":
		idx := 0
		if len(spl) > 1 {
			n, err := strconv.ParseInt(spl[1], 0, 0)
			if err != nil {
				idx = 0
			} else {
				idx = int(n)
			}
		}

		if idx > len(s.Messages) {
			idx = 0
		}

		return fmt.Sprintf("%s\r\nok\r\n", s.Messages[idx]), nil
	default:
	}

	return "ok\r\n", nil
}
