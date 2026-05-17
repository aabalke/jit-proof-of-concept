#include "funcdata.h"

//jit_amd64.s
TEXT ·callJIT(SB), 0, $56-8
    NO_LOCAL_POINTERS
    MOVQ code+0(FP), AX
    JMP AX
gocall:
    MOVQ R8, 8(SP)
    MOVQ R9, 24(SP)
    MOVQ SI, 40(SP)
    CALL CX
    MOVQ 8(SP), R8
    MOVQ 24(SP), R9
    MOVQ 40(SP), SI
    JMP (SP)

// Helper: func callJITImplAddr() uintptr
// Returns the address of the ABI0 implementation symbol.
TEXT ·callJITImplAddr(SB), 0, $0-8
    NO_LOCAL_POINTERS
    MOVQ $·callJIT(SB), AX  // address of ABI0 impl, not trampoline
    MOVQ AX, ret+0(FP)
    RET
