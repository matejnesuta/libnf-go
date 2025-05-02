package filter

import (
	LnfErr "libnf/api/errors"
	LnfRec "libnf/api/record"
	"libnf/internal"
)

type Filter struct {
	allocated bool
	ptr       uintptr
	repr      string
}

func (f Filter) String() string {
	return f.repr
}

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

func (f *Filter) Free() error {
	if !f.allocated {
		return LnfErr.ErrFilterNotInit
	}
	internal.Filter_free(f.ptr)
	f.allocated = false
	f.repr = ""
	return nil
}

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
