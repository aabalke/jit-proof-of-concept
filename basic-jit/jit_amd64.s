// jit_amd64.s
TEXT ·callJIT(SB), 0, $0-8
    MOVQ code+0(FP), AX
    JMP AX
