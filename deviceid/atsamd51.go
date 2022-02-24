//go:build atsamd51
// +build atsamd51

package main

import (
	"unsafe"
)

func deviceid() ([]uint32, error) {
	buf := [4]uint32{}

	//lint:ignore
	buf[0] = *(*uint32)(unsafe.Pointer(uintptr(0x008061FC)))
	buf[1] = *(*uint32)(unsafe.Pointer(uintptr(0x00806010)))
	buf[2] = *(*uint32)(unsafe.Pointer(uintptr(0x00806014)))
	buf[3] = *(*uint32)(unsafe.Pointer(uintptr(0x00806018)))

	return buf[:], nil
}
