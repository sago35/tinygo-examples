package main

import (
	"bufio"
	"fmt"
	"time"

	"github.com/sago35/tinygo-examples/wioterminal/initialize"
	"tinygo.org/x/drivers/net/http"
)

// You can override the setting with the init() in another source code.
// func init() {
//    ssid = "your-ssid"
//    password = "your-password"
//    url = "http://tinygo.org/"
// }

var (
	ssid     string
	password string
	url      = "http://tinygo.org/"
)

func main() {
	err := run()
	for err != nil {
		fmt.Printf("error: %s\r\n", err.Error())
		time.Sleep(5 * time.Second)
	}
}

func run() error {
	err := initialize.Wifi(ssid, password)
	if err != nil {
		return err
	}

	// You can send and receive cookies in the following way
	// 	import "tinygo.org/x/drivers/net/http/cookiejar"
	// 	jar, err := cookiejar.New(nil)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	client := &http.Client{Jar: jar}
	// 	http.DefaultClient = client

	for {
		// Various examples are as follows
		//
		// -- Get
		// 	resp, err := http.Get(url)
		//
		// -- Post
		// 	body := `cnt=12`
		// 	resp, err = http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(body))
		//
		// -- Post with JSON
		// 	body := `{"msg": "hello"}`
		// 	resp, err := http.Post(url, "application/json", strings.NewReader(body))

		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			fmt.Printf("%s\r\n", scanner.Text())
		}
		resp.Body.Close()

		time.Sleep(10 * time.Second)
	}
}
