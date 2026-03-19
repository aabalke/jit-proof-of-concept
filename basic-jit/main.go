// main.go
package main

import (
	"fmt"
	"unsafe"
)

func main() {

	c := uint64(0)
	addr := uint64(uintptr(unsafe.Pointer(&c)))

	code := []byte{
		// MOVABSQ &c, RCX  — load the address of c into RCX
		0x48, 0xB9,
		byte(addr),
		byte(addr >> 8),
		byte(addr >> 16),
		byte(addr >> 24),
		byte(addr >> 32),
		byte(addr >> 40),
		byte(addr >> 48),
		byte(addr >> 56),

		// MOVL $0xDEADBEEF, [RCX]  — move imm32 into memory at addr in RCX
		0xC7, 0x01,
		0xEF, 0xBE, 0xAD, 0xDE, // 0xDEADBEEF in little-endian

		// RET
		0xC3,
	}

	executable, err := mmapExecutable(len(code))
	if err != nil {
		panic(err)
	}

	defer munmapExecutable(executable)

	copy(executable, code)

	callJIT(&executable[0])
	fmt.Printf("C %08X\n", c)
}

// asm stub
func callJIT(code *byte)
