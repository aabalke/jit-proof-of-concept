#include "funcdata.h"

TEXT ·callJIT(SB), 0, $32-16
    NO_LOCAL_POINTERS
    MOVD code+0(FP), R0
    JMP (R0)
gocall:
    PCALIGN $16
    WORD $0xDEADBE00
    MOVD R30, 24(RSP)    // Save LR into Go local frame
    MOVD R4,  16(RSP)    // Save R4 into Go local frame
    CALL R2
    MOVD 24(RSP), R30    // Restore LR
    MOVD 16(RSP), R4     // Restore R4
    JMP (R4)
//cleanup:
//    PCALIGN $16
//    WORD $0xDEADBE01
//    ADD $32+16, RSP, RSP
//    RET

TEXT ·callJITImplAddr(SB), 0, $0-16
    NO_LOCAL_POINTERS
    MOVD $·callJIT(SB), R0  // address of ABI0 impl, not trampoline
    MOVD R0, ret+0(FP)
    RET
