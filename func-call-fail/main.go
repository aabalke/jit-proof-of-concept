// main.go
package main

import (
	"reflect"
	"runtime"
)

func main() {
	a := funcAddr(goFunction)

	code := []byte{
		// MOVABSQ a, RAX
		0x48, 0xB8,
		byte(a),
		byte(a >> 8),
		byte(a >> 16),
		byte(a >> 24),
		byte(a >> 32),
		byte(a >> 40),
		byte(a >> 48),
		byte(a >> 56),
		// CALL AX
		0xff, 0xd0,
		// RET
		0xc3,
	}

	executable, err := mmapExecutable(len(code))
	if err != nil {
		panic(err)
	}

	defer munmapExecutable(executable)

	copy(executable, code)
	callJIT(&executable[0])
}

func goFunction() {
	println("called from jit code")
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
