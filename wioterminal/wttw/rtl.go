package main

import (
	"fmt"
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

func makeHTTPRequest() string {

	var err error
	if conn != nil {
		conn.Close()
	}

	// make TCP connection
	ip := net.ParseIP(server)
	raddr := &net.TCPAddr{IP: ip, Port: 8080}
	laddr := &net.TCPAddr{Port: 8080}

	fmt.Printf("makeHTTPRequest\r\n")
	conn, err = net.DialTCP("tcp", laddr, raddr)
	for ; err != nil; conn, err = net.DialTCP("tcp", laddr, raddr) {
		time.Sleep(5 * time.Second)
	}
	println("Connected!\r")

	print("Sending HTTP request...")
	fmt.Fprintln(conn, "GET / HTTP/1.1")
	fmt.Fprintln(conn, "Host:", server)
	fmt.Fprintln(conn, "User-Agent: TinyGo/0.14.0")
	fmt.Fprintln(conn, "Connection: close")
	fmt.Fprintln(conn)
	println("Sent!\r\n\r")

	if conn == nil {
		return ""
	}

	var n int
	var ret string
	for n, err = conn.Read(buf[:]); n > 0; n, err = conn.Read(buf[:]) {
		if err != nil {
			println("Read error: " + err.Error())
		} else {
			ret = string(buf[0:n])
			print(ret)
		}
	}
	if err != nil {
		println("Read error2: " + err.Error())
	}

	err = conn.Close()
	if err != nil {
		println("Read error3: " + err.Error())
	}
	//rtl8720dn.NextSocketCh()
	return ret
}
