//go:build atsamd21
// +build atsamd21

package main

import (
	"unsafe"
)

func deviceid() ([]uint32, error) {
	buf := [4]uint32{}

	//lint:ignore
	buf[0] = *(*uint32)(unsafe.Pointer(uintptr(0x0080A00C)))
	buf[1] = *(*uint32)(unsafe.Pointer(uintptr(0x0080A040)))
	buf[2] = *(*uint32)(unsafe.Pointer(uintptr(0x0080A044)))
	buf[3] = *(*uint32)(unsafe.Pointer(uintptr(0x0080A048)))

	return buf[:], nil
}
