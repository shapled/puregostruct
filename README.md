## install

```bash
go get github.com/shapled/puregostruct@latest
```

## example

```golang
import (
	"testing"
	"unsafe"

	"github.com/shapled/puregostruct"
)

type size_t uintptr

func main() {
	var libc truct {
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
		panic(err)
	}

	if libc.Malloc == nil || libc.Free == nil || libc.Puts == nil {
		panic("libc functions not loaded")
	}

	vptr := libc.Malloc(16)
	if vptr == nil {
		panic("malloc failed")
	}

	libc.Free(vptr)

	libc.Puts("Hello, World!")
}
```
