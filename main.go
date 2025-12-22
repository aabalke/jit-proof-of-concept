package main

import (
	"bytes"
	"log"
	"reflect"
	"runtime"
	"unsafe"
)

func main() {

    // identical to tests

    runJIT(func() {
        println("called from JIT")
    })

    runJIT(func() {
        println("called from JIT, with garbage collection")
        runtime.GC()
    })

    runJIT(func() {
        println("called from JIT, with recursion (stack growth)")
        recursive()
    })
}

var i = 1 << 16
func recursive() {
    if i > 0 {
        i--
        recursive()
    }
}

func runJIT(f any) {

	a := funcAddr(f)

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
		log.Panicf("mmap: %v", err)
	}
	copy(executable, code)
	callJIT(&executable[0])

    munmapExecutable(executable)
}

func funcAddr(f any) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("funcAddr: not a func")
	}
	return v.Pointer()

    // replaces:
    //	type emptyInterface struct {
    //		typ   uintptr
    //		value *uintptr
    //	}
    //	e := (*emptyInterface)(unsafe.Pointer(&f))
    //	return *e.value
}

func getCallAddr() uintptr {

	impl := callJITImplAddr()

    // most offsets seem to be between 30 - 40
    b := unsafe.Slice((*byte)(unsafe.Pointer(impl)), 0x60)

    // equal to call cx
    p := []byte{0xFF, 0xD1}

    // get index of CALL CX
    offset := bytes.Index(b, p)

    // replaces j := funcAddr(callJIT) + 36
	j := impl + uintptr(offset)

    return j
}

// asm stubs
func callJIT(code *byte)
func callJITImplAddr() uintptr
