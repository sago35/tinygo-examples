package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	err := run()
	for err != nil {
		log.Fatal(err)
		time.Sleep(10 * time.Second)
	}
}

func run() error {
	buf, err := deviceid()
	if err != nil {
		return err
	}

	for {
		for i := range buf {
			if i == 0 {
				fmt.Printf("%08X", buf[i])
			} else {
				fmt.Printf(" %08X", buf[i])
			}
		}
		fmt.Printf("\r\n")
		time.Sleep(5 * time.Second)
	}
}
