package main

import (
	"fmt"
	"jit-test/cache"
	"reflect"
	"runtime"

	"unsafe"
)

var d uint64

//go:nosplit
func goFunction() {
	println("called from jit code")
	runtime.GC()
}

var callerPtr = uint64(getTaggedLabelAddr(0x0))

func main() {

	c := uint64(0)
	cPtr := uint64(uintptr(unsafe.Pointer(&c)))

	asm, err := New(128)
	if err != nil {
		panic(err)
	}

	asm.CallFunc(goFunction)
	asm.Mov64(R00, cPtr)
	asm.Movz(R01, 0xDEAD, 0, false)
	asm.Movk(R01, 0xBEEF, 1, false)
	asm.StrImm(R01, R00, 0, SIZE_DWRD, false, true)

	asm.Exit()

	CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

	asm.Release()

	println("completed")
	fmt.Printf("C %08X\n", c)
}

func (a *Assembler) CallFunc(f any) {
	funcPtr := uint64(funcAddr(f))
	offset := 4 * 2

	a.Mov64(R02, funcPtr)
	a.Mov64(R03, callerPtr)
	a.addInst(getADR(R04, int32(offset)))
	a.addInst(getBR(R03))
}

func (a *Assembler) Exit() {

	// this amount needs to match the amount in callJIT asm text header
	a.ADDImm(RSP, RSP, (32 + 16), false, false, true)
	a.Ret()

	if err := a.Error(); err != nil {
		panic(err)
	}

	cache.ClearICache(a.Buf)
}

func funcAddr(f any) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("funcAddr: not a func")
	}
	return v.Pointer()
}

func getADR(rd Reg, imm int32) uint32 {

	if imm >= 1 << 20 {
		panic("adr with imm >= 1 << 20")
	}

    u := uint32(imm)
    immlo := (u & 3) << 29
    immhi := (u >> 2) << 5
    return (1 << 28) | immlo | immhi | uint32(rd)
}

func getBLR(rn Reg) uint32 {
	return (uint32(0b1101_0110_0011_1111) << 16) | uint32(rn << 5)
}

func getBR(rn Reg) uint32 {
	return (uint32(0b1101_0110_0001_1111) << 16) | uint32(rn << 5)
}
