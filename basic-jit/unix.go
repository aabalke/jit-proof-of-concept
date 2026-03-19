// unix.go
//go:build linux

package main

import (
	"syscall"
	"unsafe"
)

func mmapExecutable(length int) ([]byte, error) {
	const (
		addr  = 0
		prot  = syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
		flags = syscall.MAP_PRIVATE | syscall.MAP_ANON
		fd    = 0
		off   = 0
	)

	ptr, _, err := syscall.Syscall6(
		syscall.SYS_MMAP,
		addr, uintptr(length), prot, flags, fd, off)
	if err != 0 {
		return nil, err
	}

	// Build a Go slice backed by the allocated memory
	slice := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length)
	return slice, nil
}

func munmapExecutable(_ []byte) error {
	return nil
}
