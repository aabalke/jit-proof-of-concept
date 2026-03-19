// main.go
package main

import (
	"bytes"
	"reflect"
	"runtime"
	"unsafe"
)

func main() {

	a := funcAddr(goFunction)
	j := getCallAddr()

	code := []byte{
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
		// ADDQ $framesize, SP
		0x48, 0x83, 0xc4, (8 + 8),
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
