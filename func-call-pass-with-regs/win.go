//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procVirtualAlloc = kernel32.NewProc("VirtualAlloc")
	procVirtualFree  = kernel32.NewProc("VirtualFree")
)

const (
	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
	MEM_RELEASE = 0x8000

	PAGE_EXECUTE_READWRITE = 0x40
)

func mmapExecutable(length int) ([]byte, error) {
	ptr, _, err := procVirtualAlloc.Call(
		0,
		uintptr(length),
		MEM_COMMIT|MEM_RESERVE,
		PAGE_EXECUTE_READWRITE,
	)
	if ptr == 0 {
		return nil, fmt.Errorf("VirtualAlloc failed: %w", err)
	}

	// Build a Go slice backed by the allocated memory
	slice := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length)
	return slice, nil
}

func munmapExecutable(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	addr := uintptr(unsafe.Pointer(&b[0]))
	r, _, err := procVirtualFree.Call(
		addr,
		0,
		MEM_RELEASE,
	)
	if r == 0 {
		return fmt.Errorf("VirtualFree failed: %w", err)
	}
	return nil
}
