package main

import (
	"runtime"
	"testing"
)

var called bool

func TestCall(t *testing.T) {

    // vars have to be global

    runJIT(func ()  {
        called = true
    })

    if !called {
        t.Error("Calling Func did not call inputted function")
    }

    println("completed call test")
}

func TestRecursion(t *testing.T) {
    runJIT(recursive)

    println("completed recursion test, forced stack growth")
}

func TestGC(t *testing.T) {

    runJIT(func() {
        runtime.GC()
    })

    println("completed garbage collecting test")
}
