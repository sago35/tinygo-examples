//go:build stm32f405
// +build stm32f405

package main

import (
	"unsafe"
)

func deviceid() ([]uint32, error) {
	buf := [3]uint32{}

	//lint:ignore
	buf[0] = *(*uint32)(unsafe.Pointer(uintptr(0x1FFF7A10)))
	buf[1] = *(*uint32)(unsafe.Pointer(uintptr(0x1FFF7A14)))
	buf[2] = *(*uint32)(unsafe.Pointer(uintptr(0x1FFF7A18)))

	return buf[:], nil
}
