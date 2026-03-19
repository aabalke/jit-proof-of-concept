# Go Jit Compiler Proof of Concept

A proof of concept JIT Compiler in Go meant to be supplemental to [this article](https://aaronbalke.com/posts/gojit/).
The main purpose is to illustrate the problem of calling Go functions from within JIT blocks.

This works on x86/amd64 instructions set systems, is made for Go version 1.26; however, it should work on versions >= 1.17.

This proof of concept has been developed into the [aabalke/gojit](https://github.com/aabalke/gojit) compiler.
The main usage is inside the [aabalke/guac](https://github.com/aabalke/guac) emulator, a GB, GBA, and NDS emulator.

Based on [Calling Go funcs from asm and JITed code](https://www.quasilyte.dev/blog/post/call-go-from-jit/) by Iskander Sharipov.

## Major changes
- Works with modern ABI. ABI0 and ABIInternal wrappers handled.
- Has Mmap functions for building in Linux and Windows
- Has offset calculated at runtime, since OS and Go version will cause the offset to move slightly based on the instructions used.
