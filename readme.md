# Go Jit Proof of Concept

This is a basic jit implimentation in golang, with the ability to call go functions from within the jit assembly code.
This is based on [Calling Go funcs from asm and JITed code](https://www.quasilyte.dev/blog/post/call-go-from-jit/) by Iskander Sharipov, however Sharipov's version only works on Go version <1.16. This version works on Go version 1.17+ by handling the new ABI Internal -> ABI0 changes.

An indepth article is available [here](https://aaronbalke.com/posts/calling-go-functions-from-jit-code/).

## Other Notes
- Has Mmap functions for building in linux and windows
- Has offset calculated at runtime, since OS and Go version will cause the offset to move slightly based on the instructions used.
