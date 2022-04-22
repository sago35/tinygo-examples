package main

import (
	"fmt"
	"strings"
	"time"
)

var message = "hello world"
var url = "http://192.168.1.102/"

func main() {
	err := _init()
	for err != nil {
		fmt.Printf("error : %w\r\n", err)
		time.Sleep(10 * time.Second)
	}

	body := fmt.Sprintf(`{"message": "%s"}`, message)
	post(url, "application/json", strings.NewReader(body))
}
