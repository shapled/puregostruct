package puregostruct_test

import (
	"testing"
	"unsafe"

	"github.com/shapled/puregostruct"
)

type size_t uintptr

type Libc struct {
	Malloc func(size_t) unsafe.Pointer `purego:"malloc"`
	Free   func(unsafe.Pointer)        `purego:"free"`
	Puts   func(string)                `purego:"puts"`
}

func TestLibc(t *testing.T) {
	var libc Libc
	if err := puregostruct.LoadLibrary(&libc,
		"/usr/lib/libSystem.B.dylib", // darwin
		"libc.so.7",                  // freebsd
		"libc.so.6",                  // linux
		"libc.so",                    // netbsd
		"ucrtbase.dll",               // windows
	); err != nil {
		t.Fatal(err)
	}

	if libc.Malloc == nil || libc.Free == nil || libc.Puts == nil {
		t.Fatal("libc functions not loaded")
	}

	vptr := libc.Malloc(16)
	if vptr == nil {
		t.Fatal("malloc failed")
	}

	libc.Free(vptr)

	libc.Puts("Hello, World!")
}
