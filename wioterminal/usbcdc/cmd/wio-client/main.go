package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := flag.String("port", "", "port")
	flag.Parse()

	p, err := serial.Open(*port, &serial.Mode{BaudRate: 115200})
	if err != nil {
		return err
	}

	ch := make(chan string)

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			ch <- "led"
		}
	}()

	go func() {
		for {
			time.Sleep(777 * time.Millisecond)
			ch <- "lcd"
		}
	}()

	go func() {
		i := 0
		for {
			time.Sleep(1 * time.Second)
			ch <- fmt.Sprintf("msg %d", i)
			i = 1 - i
		}
	}()

	scanner := bufio.NewScanner(p)
	for cmd := range ch {
		fmt.Fprintf(p, "%s\r\n", cmd)
		for {
			scanner.Scan()
			ret := scanner.Text()
			if ret == "ok" {
				break
			} else {
				fmt.Printf("%s\n", ret)
			}
		}
	}

	return nil
}
