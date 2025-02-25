package libnf

import (
	"libnf/internal"
	"unsafe"
)

type MemHeap struct {
	ptr       uintptr
	allocated bool
}

type MemHeapCursor struct {
	ptr uintptr
}

const (
	FastAggrNone  int = internal.FAST_AGGR_NONE
	FastAggrBasic int = internal.FAST_AGGR_BASIC
	FastAggrAll   int = internal.FAST_AGGR_ALL
)

const (
	SortAsc  int = internal.SORT_ASC
	SortDesc int = internal.SORT_DESC
)

const (
	AggrKey int = internal.AGGR_KEY
	AggrSum int = internal.AGGR_SUM
	AggrMin int = internal.AGGR_MIN
	AggrMax int = internal.AGGR_MAX
	AggrOr  int = internal.AGGR_OR
)

func NewMemHeap() (MemHeap, error) {
	memHeap := MemHeap{}
	status := internal.Mem_init(&memHeap.ptr)
	if status == internal.ERR_NOMEM {
		return memHeap, ErrNoMem
	} else if status == internal.ERR_OTHER {
		return memHeap, ErrOther
	}

	memHeap.allocated = true
	return memHeap, nil
}

func (m *MemHeap) Free() error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	internal.Mem_free(m.ptr)
	m.allocated = false
	return nil
}

func (m *MemHeap) Clear() error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	internal.Mem_clean(m.ptr)
	return nil
}

// int 	lnf_mem_first_c (lnf_mem_t *lnf_mem, lnf_mem_cursor_t **cursor)
//
//	Set the cursor position to the first record.
func (m *MemHeap) FirstRecordPosition() (MemHeapCursor, error) {
	if !m.allocated {
		return MemHeapCursor{}, ErrMemHeapNotAllocated
	}
	cursor := MemHeapCursor{}
	status := internal.Mem_first_c(m.ptr, &cursor.ptr)
	if status == internal.EOF {
		return cursor, ErrFileEof
	} else if status == internal.ERR_NOMEM {
		return cursor, ErrNoMem
	}
	return cursor, nil
}

// int 	lnf_mem_next_c (lnf_mem_t *lnf_mem, lnf_mem_cursor_t **cursor)
//
//	Set the cursor position to the next record.
func (m *MemHeap) NextRecordPosition(c *MemHeapCursor) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	status := internal.Mem_next_c(m.ptr, &c.ptr)
	if status == internal.EOF {
		return ErrFileEof
	} else if status == internal.ERR_NOMEM {
		return ErrNoMem
	}
	return nil
}

// int 	lnf_mem_read (lnf_mem_t *lnf_mem, lnf_rec_t *rec)
//
//	Read next record from memheap.
func (m *MemHeap) GetNextRecord(r *Record) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	} else if !r.allocated {
		return ErrRecordNotAllocated
	}
	status := internal.Mem_read(m.ptr, r.ptr)
	if status == internal.ERR_NOMEM {
		return ErrNoMem
	} else if status == internal.EOF {
		return ErrMemHeapEnd
	}
	return nil
}

// int 	lnf_mem_write (lnf_mem_t *lnf_mem, lnf_rec_t *rec)
//  	Write record to memheap object.

func (m *MemHeap) WriteRecord(r *Record) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	} else if !r.allocated {
		return ErrRecordNotAllocated
	}
	status := internal.Mem_write(m.ptr, r.ptr)
	if status == internal.ERR_NOMEM {
		return ErrNoMem
	} else if status == internal.ERR_OTHER {
		return ErrOther
	}
	return nil
}

// int 	lnf_mem_read_c (lnf_mem_t *lnf_mem, lnf_mem_cursor_t *cursor, lnf_rec_t *rec)
//
//	Read next record on the position given by cursor.
func (m *MemHeap) GetRecordWithCursor(c *MemHeapCursor, r *Record) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	} else if !r.allocated {
		return ErrRecordNotAllocated
	}
	status := internal.Mem_read_c(m.ptr, c.ptr, r.ptr)
	if status == internal.ERR_NOMEM {
		return ErrNoMem
	} else if status == internal.EOF {
		return ErrMemHeapEnd
	}
	return nil
}

// int 	lnf_mem_lookup_c (lnf_mem_t *lnf_mem, lnf_rec_t *rec, lnf_mem_cursor_t **cursor)
//
//	Set the cursor position to the record identified by key fields.
func (m *MemHeap) GetRecordWithKey(r *Record) (MemHeapCursor, error) {
	if !m.allocated {
		return MemHeapCursor{}, ErrMemHeapNotAllocated
	} else if !r.allocated {
		return MemHeapCursor{}, ErrRecordNotAllocated
	}
	cursor := MemHeapCursor{}
	status := internal.Mem_lookup_c(m.ptr, r.ptr, &cursor.ptr)
	if status == internal.EOF {
		return cursor, ErrMemHeapEnd
	} else if status == internal.ERR_NOMEM {
		return cursor, ErrNoMem
	}
	return cursor, nil
}

// int 	lnf_mem_merge_threads (lnf_mem_t *lnf_mem)
//
//	Merge data from multiple threads into one thread.
func (m *MemHeap) MergeThreads() error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	status := internal.Mem_merge_threads(m.ptr)
	if status == internal.ERR_NOMEM {
		return ErrNoMem
	} else if status == internal.EOF {
		return ErrMemHeapEnd
	}
	return nil
}

// int 	lnf_mem_fastaggr (lnf_mem_t *lnf_mem, int flags)
//
//	Set fast aggregation mode.
func (m *MemHeap) SetFastAggr(option int) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	status := internal.Mem_fastaggr(m.ptr, option)
	if status == internal.ERR_NOMEM {
		return ErrNoMem
	} else if status == internal.ERR_OTHER {
		return ErrOther
	}
	return nil
}

// int 	lnf_mem_setopt (lnf_mem_t *lnf_mem, int opt, void *data, size_t size)
//  	Set lnf_mem_t options.

func callMemHeapSetOpt(m *MemHeap, opt int, data uintptr, size int64) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	status := internal.Mem_setopt(m.ptr, opt, data, size)
	if status == internal.ERR_OTHER {
		return ErrOther
	}
	return nil
}

func (m *MemHeap) SetListMode() error {
	return callMemHeapSetOpt(m, internal.OPT_LISTMODE, uintptr(unsafe.Pointer(nil)), 0)
}

func (m *MemHeap) SetHashBuckets(num int) error {
	return callMemHeapSetOpt(m, internal.OPT_HASHBUCKETS, uintptr(unsafe.Pointer(&num)), int64(unsafe.Sizeof(num)))
}

func (m *MemHeap) EnableNfdumpCompat() error {
	return callMemHeapSetOpt(m, internal.OPT_COMP_STATSCMP, uintptr(unsafe.Pointer(nil)), 0)
}

// int 	lnf_mem_fadd (lnf_mem_t *lnf_mem, int field, int flags, int numbits, int numbits6)
//
//	Set aggregation and sort option for memheap.
func (m *MemHeap) SetAggrOptions(field int, aggrType int, sortType int, numBits int, numBits6 int) error {
	if !m.allocated {
		return ErrMemHeapNotAllocated
	}
	status := internal.Mem_fadd(m.ptr, field, aggrType+sortType, numBits, numBits6)
	if status == internal.ERR_OTHER {
		return ErrOther
	} else if status == internal.ERR_NOMEM {
		return ErrNoMem
	}
	return nil
}
