package main

import (
	"fmt"
	"strings"
	"time"

	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/rtl8720dn"
)

// IP address of the server aka "hub". Replace with your own info.
var server = "192.168.1.114"
var conn net.Conn
var buf [4096]byte

// connect to RTL8720DN
func connectToRTL8720() bool {
	for i := 0; i < 5; i++ {
		println("Connecting to wifi adaptor...")
		if rtl.Connected() {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

// connect to access point
func connectToAP() {
	println("Connecting to wifi network '" + ssid + "'")

	err := rtl.SetWifiMode(rtl8720dn.WifiModeClient)
	if err != nil {
		failMessage(err.Error())
	}

	for {
		err = rtl.ConnectToAP(ssid, pass, 40)
		if err != nil {
			fmt.Printf("%s\r\n", err.Error())
			//failMessage(err.Error())
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}

	println("Connected.")
	ip, err := rtl.GetClientIP()
	if err != nil {
		failMessage(err.Error())
	}

	println(ip)
}

func failMessage(msg string) {
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}

func httpGet(url string) ([]byte, error) {
	var err error
	if conn != nil {
		conn.Close()
	}

	// make TCP connection
	x := strings.SplitN(url, "://", 2)
	if 1 == len(x) {
		return nil, fmt.Errorf("url error: %q", url)
	}
	x = strings.SplitN(x[1], "/", 2)
	address := x[0]
	path := "/"
	if 1 < len(x) {
		path += x[1]
	}

	raddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\r\n", raddr)
	laddr := &net.TCPAddr{Port: 8080}

	conn, err = net.DialTCP("tcp", laddr, raddr)
	for ; err != nil; conn, err = net.DialTCP("tcp", laddr, raddr) {
		time.Sleep(5 * time.Second)
	}

	fmt.Fprintln(conn, fmt.Sprintf("GET %s HTTP/1.1", path))
	fmt.Fprintln(conn, "Host:", server)
	fmt.Fprintln(conn, "User-Agent: TinyGo/0.14.0")
	fmt.Fprintln(conn, "Connection: close")
	fmt.Fprintln(conn)

	if conn == nil {
		return buf[:0], nil
	}

	var n int
	var end int
	rbuf := buf[:]

	n, err = conn.Read(rbuf)
	for n > 0 {
		if err != nil {
			println("Read error: " + err.Error())
		} else {
			end += n
		}
		n, err = conn.Read(rbuf)
	}
	if err != nil {
		println("Read error2: " + err.Error())
	}

	err = conn.Close()
	if err != nil {
		println("Read error3: " + err.Error())
	}
	return buf[:end], nil
}
