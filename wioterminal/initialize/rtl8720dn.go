//go:build wioterminal
// +build wioterminal

package initialize

import (
	"device/sam"
	"machine"
	"runtime"
	"runtime/interrupt"
	"time"

	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/net/http"
	"tinygo.org/x/drivers/rtl8720dn"
)

var (
	rtl       *rtl8720dn.Driver
	connected bool
	uart      UARTx
	debug     bool
	buf       [0x1000]byte
)

func handleInterrupt(interrupt.Interrupt) {
	// should reset IRQ
	uart.Receive(byte((uart.Bus.DATA.Get() & 0xFF)))
	uart.Bus.INTFLAG.SetBits(sam.SERCOM_USART_INT_INTFLAG_RXC)
}

// SetupRTL8720DN sets up the RTL8270DN for use.
func SetupRTL8720DN() (*rtl8720dn.Driver, error) {
	rtl := rtl8720dn.New(machine.UART3, machine.PB24, machine.PC24, machine.RTL8720D_CHIP_PU)

	if debug {
		waitSerial()
	}

	rtl.Debug(debug)
	rtl.Configure()

	connected = true
	return rtl, nil
}

// Wifi sets up the RTL8720DN and connects it to Wi-Fi.
func Wifi(ssid, pass string, timeout time.Duration) (*rtl8720dn.Driver, error) {
	_, err := SetupRTL8720DN()
	if err != nil {
		return nil, err
	}

	err = rtl.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil {
		return rtl, err
	}

	net.UseDriver(rtl)
	http.UseDriver(rtl)
	http.SetBuf(buf[:])

	// NTP
	t, err := GetCurrentTime()
	if err != nil {
		return nil, err
	}
	runtime.AdjustTimeOffset(-1 * int64(time.Since(t)))

	return rtl, nil
}

func Device() *rtl8720dn.Driver {
	return rtl
}

func Connected() bool {
	return connected
}

func IP() rtl8720dn.IPAddress {
	ip, _, _, _ := rtl.GetIP()
	return ip
}

func Subnet() rtl8720dn.IPAddress {
	_, subnet, _, _ := rtl.GetIP()
	return subnet
}

func Gateway() rtl8720dn.IPAddress {
	_, _, gateway, _ := rtl.GetIP()
	return gateway
}

func SetRootCA(s *string) {
	rtl.SetRootCA(s)
}

// Wait for user to open serial console
func waitSerial() {
	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}
}

type UARTx struct {
	*machine.UART
}

func (u UARTx) Read(p []byte) (n int, err error) {
	if u.Buffered() == 0 {
		time.Sleep(1 * time.Millisecond)
		return 0, nil
	}
	return u.UART.Read(p)
}

// Debug sets the debug mode.
func Debug(b bool) {
	debug = b
}
