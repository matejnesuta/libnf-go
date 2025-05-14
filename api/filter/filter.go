package filter

import (
	LnfErr "github.com/matejnesuta/libnf-go/api/errors"
	LnfRec "github.com/matejnesuta/libnf-go/api/record"
	"github.com/matejnesuta/libnf-go/internal"
)

// Filter represents a compiled flow record filter.
//
// It can be used to match flow records using a flexible filtering expression.
// Internally, it uses the newer libnf filter engine (v2), which is thread-safe,
// leak-free, and more extensible compared to the legacy nfdump filter code.
type Filter struct {
	allocated bool
	ptr       uintptr
	repr      string
}

// String returns the original filter expression used to initialize the Filter.
func (f Filter) String() string {
	return f.repr
}

// Init compiles the provided filter expression and initializes the Filter.
//
// This uses the new libnf v2 filtering engine, which supports multithreading
// and proper memory cleanup. The filter expression syntax is similar to that
// used in nfdump, though not all legacy features may be available.
//
// Returns an error if memory allocation fails, the filter expression is invalid,
// or if the filter is already initialized.
func (f *Filter) Init(expression string) error {
	if f.allocated {
		return LnfErr.ErrFilterAlreadyInit
	}

	status := internal.Filter_init_v2(&f.ptr, expression)
	if status == internal.ERR_NOMEM {
		return LnfErr.ErrNoMem
	} else if status == internal.ERR_FILTER {
		return LnfErr.ErrFilter
	} else if status == internal.ERR_OTHER_MSG {
		return LnfErr.ErrOtherMsg
	}
	f.repr = expression
	f.allocated = true
	return nil
}

// Free releases the resources allocated for the Filter.
//
// After calling Free, the Filter must be reinitialized before use.
// If the Filter was not initialized, this returns an error.
func (f *Filter) Free() error {
	if !f.allocated {
		return LnfErr.ErrFilterNotInit
	}
	internal.Filter_free(f.ptr)
	f.allocated = false
	f.repr = ""
	return nil
}

// Match checks whether the given flow record satisfies the filter criteria.
//
// Returns true if the record matches the filter expression, false otherwise.
// Returns an error if the filter is not initialized or the record is not allocated.
func (f *Filter) Match(r LnfRec.Record) (bool, error) {
	if !f.allocated {
		return false, LnfErr.ErrFilterNotInit
	} else if !r.Allocated() {
		return false, LnfErr.ErrRecordNotAllocated
	}
	status := internal.Filter_match(f.ptr, r.GetPtr())
	if status == 1 {
		return true, nil
	}
	return false, nil
}
