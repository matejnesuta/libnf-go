package libnf

import (
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
		return ErrFilterAlreadyInit
	}

	status := internal.Filter_init_v2(&f.ptr, expression)
	if status == internal.ERR_NOMEM {
		return ErrNoMem
	} else if status == internal.ERR_FILTER {
		return ErrFilter
	} else if status == internal.ERR_OTHER_MSG {
		return ErrOtherMsg
	}
	f.repr = expression
	f.allocated = true
	return nil
}

func (f *Filter) Free() error {
	if !f.allocated {
		return ErrFilterNotInit
	}
	internal.Filter_free(f.ptr)
	f.allocated = false
	return nil
}

func (f Filter) Match(r Record) (bool, error) {
	if !f.allocated {
		return false, ErrFilterNotInit
	} else if !r.allocated {
		return false, ErrRecordNotAllocated
	}
	status := internal.Filter_match(f.ptr, r.ptr)
	if status == 1 {
		return true, nil
	}
	return false, nil
}
