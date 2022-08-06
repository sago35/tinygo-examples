package initialize

import (
	"errors"
	"fmt"
	"time"

	"tinygo.org/x/drivers/net"
)

const (
	ntpPacketSize = 48
	ntpHost       = "129.6.15.29"
)

var (
	b = make([]byte, ntpPacketSize)
)

func GetCurrentTime() (time.Time, error) {
	hip := net.ParseIP(ntpHost)
	raddr := &net.UDPAddr{IP: hip, Port: 123}
	laddr := &net.UDPAddr{Port: 2390}
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return time.Time{}, err
	}
	defer conn.Close()

	return getCurrentTime(conn)
}

func getCurrentTime(conn *net.UDPSerialConn) (time.Time, error) {
	if err := sendNTPpacket(conn); err != nil {
		return time.Time{}, err
	}
	clearBuffer()
	for now := time.Now(); time.Since(now) < time.Second; {
		time.Sleep(5 * time.Millisecond)
		if n, err := conn.Read(b); err != nil {
			return time.Time{}, fmt.Errorf("error reading UDP packet: %w", err)
		} else if n == 0 {
			continue // no packet received yet
		} else if n != ntpPacketSize {
			if n != ntpPacketSize {
				return time.Time{}, fmt.Errorf("expected NTP packet size of %d: %d", ntpPacketSize, n)
			}
		}
		return parseNTPpacket(), nil
	}
	return time.Time{}, errors.New("no packet received after 1 second")
}

func sendNTPpacket(conn *net.UDPSerialConn) error {
	clearBuffer()
	b[0] = 0b11100011 // LI, Version, Mode
	b[1] = 0          // Stratum, or type of clock
	b[2] = 6          // Polling Interval
	b[3] = 0xEC       // Peer Clock Precision
	// 8 bytes of zero for Root Delay & Root Dispersion
	b[12] = 49
	b[13] = 0x4E
	b[14] = 49
	b[15] = 52
	if _, err := conn.Write(b); err != nil {
		return err
	}
	return nil
}

func parseNTPpacket() time.Time {
	// the timestamp starts at byte 40 of the received packet and is four bytes,
	// this is NTP time (seconds since Jan 1 1900):
	t := uint32(b[40])<<24 | uint32(b[41])<<16 | uint32(b[42])<<8 | uint32(b[43])
	const seventyYears = 2208988800
	return time.Unix(int64(t-seventyYears), 0)
}

func clearBuffer() {
	for i := range b {
		b[i] = 0
	}
}
