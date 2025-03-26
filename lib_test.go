package puregostruct_test

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/shapled/puregostruct"
)

type size_t uintptr

func TestLoadLibrary(t *testing.T) {
	var libc struct {
		Malloc func(size_t) unsafe.Pointer `purego:"malloc"`
		Free   func(unsafe.Pointer)        `purego:"free"`
		Puts   func(string)                `purego:"puts"`
	}
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

func TestLoadLibraryPanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expected panic")
		}
		msg := fmt.Sprintf("%v", err)
		if !strings.HasPrefix(msg, "failed to register field `NoSuchFunc` with purego tag `no_such_func`") {
			t.Fatalf("wrong error message: %s", msg)
		}
	}()

	var libc struct {
		NoSuchFunc func(size_t) unsafe.Pointer `purego:"no_such_func"`
	}
	puregostruct.LoadLibrary(&libc,
		"/usr/lib/libSystem.B.dylib", // darwin
		"libc.so.7",                  // freebsd
		"libc.so.6",                  // linux
		"libc.so",                    // netbsd
		"ucrtbase.dll",               // windows
	)
}
