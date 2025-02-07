package internal

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lnf
// #include "libnf.h"
import "C"
import (
	"unsafe"
)

func Error() string {
	buf := make([]byte, MAX_STRING)                                    // Allocate buffer in Go
	C.lnf_error((*C.char)(unsafe.Pointer(&buf[0])), C.int(MAX_STRING)) // Call SWIG-wrapped C function
	return string(buf)                                                 // Convert C buffer to Go string
}
