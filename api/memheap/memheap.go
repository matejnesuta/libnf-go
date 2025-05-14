// Memheap is the set of methods that allows to aggregate and sort
// records data. This is usually done in five steps:
//
// 1. Create new memory heap using the NewMemHeap function.
//
// 2. Set key, aggregation and sort key via SetAggrOptions function.
//
// 3. Lock to an OS thread using the runtime.LockOSThread function and fill the internal structure with input records via WriteRecord.
//
// 4. Read aggregated and sorted result via GetNextRecord.
//
// 5. Release Memheap structure and all relevant resources using the Free function.
package memheap

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lnf
// #include "libnf.h"
import "C"

import (
	"unsafe"

	"github.com/matejnesuta/libnf-go/api/errors"
	"github.com/matejnesuta/libnf-go/api/record"
	"github.com/matejnesuta/libnf-go/internal"
)

// MemHeap represents a memory heap used for storing flow records.
type MemHeap struct {
	ptr       uintptr
	allocated bool
}

// MemHeapCursor represents a cursor used to navigate through records in the memory heap.
type MemHeapCursor struct {
	ptr uintptr
}

const (
	// No aggregation.
	FastAggrNone int = internal.FAST_AGGR_NONE
	// Perform aggregation on items fields.First, fields.Last, fields.Doctets, fields.Pkts.
	FastAggrBasic int = internal.FAST_AGGR_BASIC
	// Perform aggregation on all fields.
	FastAggrAll int = internal.FAST_AGGR_ALL
)

const (
	// Do not sort the result.
	SortNone int = internal.SORT_NONE
	// Sort the result in the ascending order.
	SortAsc int = internal.SORT_ASC
	// Sort the result in the descending order.
	SortDesc int = internal.SORT_DESC
)

const (
	// Default aggregation option for the field.
	AggrAuto int = internal.AGGR_AUTO
	// // Use the field as a key for aggregation.
	AggrKey int = internal.AGGR_KEY
	// Make summary of all aggregated values.
	AggrSum int = internal.AGGR_SUM
	// Find minimum value of the field.
	AggrMin int = internal.AGGR_MIN
	// Find maximum value of the field.
	AggrMax int = internal.AGGR_MAX
	// Perform OR operation on all of the values.
	AggrOr int = internal.AGGR_OR
)

// Check if the memory heap is allocated.
func (m *MemHeap) Allocated() bool {
	return m.allocated
}

// Initialize empty memheap object and allocate all necessary resources.
func NewMemHeap() (MemHeap, error) {
	memHeap := MemHeap{}
	status := internal.Mem_init(&memHeap.ptr)
	if status == internal.ERR_NOMEM {
		return memHeap, errors.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return memHeap, errors.ErrOther
	}

	memHeap.allocated = true
	return memHeap, nil
}

// Free all resources allocated for the MemHeap object.
func (m *MemHeap) Free() error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	internal.Mem_free(m.ptr)
	m.allocated = false
	return nil
}

// Clean all data in memheap. The memheap will be in the same state as it was after the initialization.
func (m *MemHeap) Clear() error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	internal.Mem_clean(m.ptr)
	return nil
}

// Set the cursor position to the first record in MemHeap.
func (m *MemHeap) FirstRecordPosition() (MemHeapCursor, error) {
	if !m.allocated {
		return MemHeapCursor{}, errors.ErrMemHeapNotAllocated
	}
	cursor := MemHeapCursor{}
	status := internal.Mem_first_c(m.ptr, &cursor.ptr)
	if status == internal.EOF {
		return cursor, errors.ErrMemHeapEnd
	} else if status == internal.ERR_NOMEM {
		return cursor, errors.ErrNoMem
	}
	return cursor, nil
}

// Set the cursor position to the next record in MemHeap.
func (m *MemHeap) NextRecordPosition(c *MemHeapCursor) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	status := internal.Mem_next_c(m.ptr, &c.ptr)
	if status == internal.EOF {
		return errors.ErrMemHeapEnd
	} else if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	}
	return nil
}

// Read next record from the MemHeap.
func (m *MemHeap) GetNextRecord(r *record.Record) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	} else if !r.Allocated() {
		return errors.ErrRecordNotAllocated
	}
	status := internal.Mem_read(m.ptr, r.GetPtr())
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.EOF {
		return errors.ErrMemHeapEnd
	}
	return nil
}

// Write record to the MemHeap object. A thread must be locked to the OS thread using runtime.LockOSThread() before calling this function.
// It is possible to call this function from multiple goroutines, but each goroutine must lock to the OS thread before calling this function and the MergeThreads function must be called at the end of each goroutine.
func (m *MemHeap) WriteRecord(r *record.Record) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	} else if !r.Allocated() {
		return errors.ErrRecordNotAllocated
	}

	status := int(C.lnf_mem_write(
		(unsafe.Pointer(m.ptr)),
		(unsafe.Pointer(r.GetPtr())),
	))
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return errors.ErrOther
	}
	return nil
}

// Read next record on the position given by cursor.
func (m *MemHeap) GetRecordWithCursor(c *MemHeapCursor, r *record.Record) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	} else if !r.Allocated() {
		return errors.ErrRecordNotAllocated
	}
	status := internal.Mem_read_c(m.ptr, c.ptr, r.GetPtr())
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.EOF {
		return errors.ErrMemHeapEnd
	}
	return nil
}

// Set the cursor position to the record identified by key fields.
func (m *MemHeap) GetRecordWithKey(r *record.Record) (MemHeapCursor, error) {
	if !m.allocated {
		return MemHeapCursor{}, errors.ErrMemHeapNotAllocated
	} else if !r.Allocated() {
		return MemHeapCursor{}, errors.ErrRecordNotAllocated
	}
	cursor := MemHeapCursor{}
	status := internal.Mem_lookup_c(m.ptr, r.GetPtr(), &cursor.ptr)
	if status == internal.EOF {
		return cursor, errors.ErrMemHeapEnd
	} else if status == internal.ERR_NOMEM {
		return cursor, errors.ErrNoMem
	}
	return cursor, nil
}

// When multiple goroutines are used to write records to the same heap, this function must be called at the end of each goroutine.
func (m *MemHeap) MergeThreads() error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	status := internal.Mem_merge_threads(m.ptr)
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.EOF {
		return errors.ErrMemHeapEnd
	}
	return nil
}

// Set fast aggregation mode.
func (m *MemHeap) SetFastAggr(option int) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	status := internal.Mem_fastaggr(m.ptr, option)
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return errors.ErrOther
	}
	return nil
}

func callMemHeapSetOpt(m *MemHeap, opt int, data uintptr, size int64) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	status := internal.Mem_setopt(m.ptr, opt, data, size)
	if status == internal.ERR_OTHER {
		return errors.ErrOther
	}
	return nil
}

// Set List mode. This is used for sorting without aggregation.
func (m *MemHeap) SetListMode() error {
	return callMemHeapSetOpt(m, internal.OPT_LISTMODE, uintptr(unsafe.Pointer(nil)), 0)
}

// Set the number of hash buckets.
func (m *MemHeap) SetHashBuckets(num int) error {
	return callMemHeapSetOpt(m, internal.OPT_HASHBUCKETS, uintptr(unsafe.Pointer(&num)), int64(unsafe.Sizeof(num)))
}

// Statistics for pair fields are counted twice in Nfdump, but only once if the pair fields are the same.
// Libnf counts the records with pair fields twice by default.
// This option switches Libnf to the Nfdump behavior.
//
// Example:
// We have input flow
// SRC           DST          PKTS BYTES
// 1.1.1.1:53 -> 2.2.2.2:53      1   20
// 3.3.3.3:80 -> 4.4.4.4:1222    3   80
//
// and we have statistics via port field
//
// In the Nfdump and Libnf with EnableNfdumpCompat option enabled
// the result will be:
// PORT    PKTS BYTES
// 53         1    20
// 80         3    80
// 1222       3    80
//
// but in Libnf without the EnableNfdumpCompat option the result will be:
// PORT    PKTS BYTES
// 53         2    40
// 80         3    80
// 1222       3    80
func (m *MemHeap) EnableNfdumpCompat() error {
	return callMemHeapSetOpt(m, internal.OPT_COMP_STATSCMP, uintptr(unsafe.Pointer(nil)), 0)
}

// SetAggrOptions configures aggregation options for a specific field in the memory heap.
//
// Parameters:
//   - field: the field ID.
//   - aggrType: the aggregation type (e.g., sum, min, max, key).
//   - sortType: the sort type to apply to the aggregation (e.g., ascending, descending).
//   - numBits: the number of bits used when working with IPv4 addresses.
//   - numBits6: the number of bits used when working with IPv6 addresses.
//
// Returns an error if the MemHeap has not been allocated or if the internal call fails.
func (m *MemHeap) SetAggrOptions(field int, aggrType int, sortType int, numBits int, numBits6 int) error {
	if !m.allocated {
		return errors.ErrMemHeapNotAllocated
	}
	status := internal.Mem_fadd(m.ptr, field, aggrType|sortType, numBits, numBits6)
	if status == internal.ERR_OTHER {
		return errors.ErrOther
	} else if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	}
	return nil
}
