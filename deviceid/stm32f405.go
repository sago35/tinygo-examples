//go:build nrf52840
// +build nrf52840

package main

import (
	"unsafe"
)

func deviceid() ([]uint32, error) {
	buf := [2]uint32{}

	//lint:ignore
	buf[0] = *(*uint32)(unsafe.Pointer(uintptr(0x10000060)))
	buf[1] = *(*uint32)(unsafe.Pointer(uintptr(0x10000064)))

	return buf[:], nil
}
