package errors

import (
	"errors"

	"github.com/matejnesuta/libnf-go/internal"
)

// Weak errors
var (
	ErrUnknownBlock  = errors.New("weak error: unknown block type")
	ErrUnknownRecord = errors.New("weak error: unknown record type")
	ErrCompat15      = errors.New("weak error: old block type supported by nfdump 1.5")
	ErrWeak          = errors.New("multiple weak errors (errors to skip)")
)

// IO and corruption errors
var (
	ErrRead    = errors.New("read error (IO)")
	ErrCorrupt = errors.New("corrupted file")
	ErrExtMapB = errors.New("extension map is too big")
	ErrExtMapM = errors.New("extension map is missing")
	ErrWrite   = errors.New("write error")
)

// File errors
var (
	ErrCannotOpenFile    = errors.New("cannot open file")
	ErrFileNotOpened     = errors.New("file is not opened")
	ErrFileAlreadyOpened = errors.New("file is already opened")
	ErrFileEof           = errors.New("end of file")
)

// Record errors
var (
	ErrRecordNotAllocated   = errors.New("record is not allocated")
	ErrUnknownFldType       = errors.New("unknown field type")
	ErrUnknownFld           = errors.New("unknown field")
	ErrMismatchingDataTypes = errors.New("expected field type does not match the actual data type")
)

// Filter errors
var (
	ErrFilter            = errors.New("cannot compile a filter")
	ErrFilterNotInit     = errors.New("filter is not initialized")
	ErrFilterAlreadyInit = errors.New("filter is already initialized")
)

// MemHeap errors
var (
	ErrMemHeapNotAllocated = errors.New("memory heap is not allocated")
	ErrMemHeapEnd          = errors.New("end of memory heap")
	ErrMemHeapEmpty        = errors.New("memory heap is empty")
)

// Other errors
var (
	ErrNotSet   = errors.New("item is not set")
	ErrNoMem    = errors.New("cannot allocate memory")
	ErrOther    = errors.New("other error")
	ErrOtherMsg = errors.New("other error with additional information")
	ErrNaN      = errors.New("attempt to divide by 0")
)

func Error() string {
	return internal.Error()
}
