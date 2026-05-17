package main

import (
	"bytes"
	"fmt"
	"unsafe"
)

// PageSize is the size of a memory page. The len argument to Alloc
// should be an integer multiple of the page size.
const PageSize = 4096

func callJIT(code uintptr)
func callJITImplAddr() uintptr

func getTaggedLabelAddr(tagIdx uint8) uintptr {
	impl := callJITImplAddr()
	bts := unsafe.Slice((*uint8)(unsafe.Pointer(impl)), 0x100)
	tagBytes := []uint8{tagIdx, 0xBE, 0xAD, 0xDE}
	offset := bytes.Index(bts, tagBytes)
	offset += 4 // past offset
	return impl + uintptr(offset)
}

func getCallAddr() uintptr {

	impl := callJITImplAddr()

	return impl + uintptr(8 * 4) + 4


	// most offsets seem to be between 30 - 40
	insts := unsafe.Slice((*uint32)(unsafe.Pointer(impl)), 0x60)

	label := getBLR(R02)
	//label := getBR(RSP)

	// get index of CALL CX (BLR R10)

	i := 0
	for ; i < len(insts); i++ {
		if insts[i] == label {
			break
		}
	}

	offset := i * 4

	fmt.Printf("Label %08X, Offset %04d\n", label, offset)

	return impl + uintptr(offset)
}
