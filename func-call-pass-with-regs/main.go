// main.go
package main

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

var c uint64

func main() {
	a := funcAddr(goFunction)
	//j := getCallAddr()
	j := getTaggedLabelAddr(0x0)

	cPtr := uint64(uintptr(unsafe.Pointer(&c)))

	v := uint64(0xDEADBEEF)

	code := []byte{
		0x49, 0xB8,
		byte(cPtr),
		byte(cPtr >> 8),
		byte(cPtr >> 16),
		byte(cPtr >> 24),
		byte(cPtr >> 32),
		byte(cPtr >> 40),
		byte(cPtr >> 48),
		byte(cPtr >> 56),
		0x49, 0xB9,
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
		byte(v >> 24),
		byte(v >> 32),
		byte(v >> 40),
		byte(v >> 48),
		byte(v >> 56),
		// MOVABSQ funcAddr(f), CX
		0x48, 0xB9,
		byte(a),
		byte(a >> 8),
		byte(a >> 16),
		byte(a >> 24),
		byte(a >> 32),
		byte(a >> 40),
		byte(a >> 48),
		byte(a >> 56),
		// MOVABSQ funcAddr(callJIT)+offset (gocall label), DI
		0x48, 0xBF,
		byte(j),
		byte(j >> 8),
		byte(j >> 16),
		byte(j >> 24),
		byte(j >> 32),
		byte(j >> 40),
		byte(j >> 48),
		byte(j >> 56),
		// LEAQ 6(PC), SI
		0x48, 0x8d, 0x35, (4 + 2), 0, 0, 0,
		// MOVQ SI, (SP)
		0x48, 0x89, 0x34, 0x24,
		// JMP DI
		0xff, 0xe7,

		// mov [r8], r9
		0x4D, 0x89, 0x08,

		// ADDQ $framesize, SP
		0x48, 0x83, 0xc4, (56 + 8),
		// RET
		0xc3,
	}

	executable, err := mmapExecutable(len(code))
	if err != nil {
		panic(err)
	}
	copy(executable, code)
	callJIT(&executable[0])

	munmapExecutable(executable)

	fmt.Printf("C %08X\n", c)
}

func goFunction() {
	println("called from jit code 1")
	runtime.GC() // the line that causes the stack functions which break the jit
}

// asm stubs
func callJIT(code *byte)

func funcAddr(f any) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("funcAddr: not a func")
	}
	return v.Pointer()
}

// asm stub
func callJITImplAddr() uintptr

func getCallAddr() uintptr {
	impl := callJITImplAddr()

	// most offsets seem to be between 30 - 40
	b := unsafe.Slice((*byte)(unsafe.Pointer(impl)), 0x60)

	// equal to call cx
	label := []byte{0xFF, 0xD1}

	// get index of CALL CX
	offset := bytes.Index(b, label)

	return impl + uintptr(offset)
}

func getTaggedLabelAddr(tagIdx uint8) uintptr {
	impl := callJITImplAddr()
	bts := unsafe.Slice((*uint8)(unsafe.Pointer(impl)), 0x60)
	//tagBytes := []uint8{tagIdx, 0xBE, 0xAD, 0xDE}
	//tagBytes := []uint8{0xDE, 0xAD, 0xBE, tagIdx}

	tagBytes := []byte{
		0x4C, 0x89, 0x44, 0x24, 0x08,
	}
	offset := bytes.Index(bts, tagBytes)
	//offset += 4 // past offset
	return impl + uintptr(offset)
}
