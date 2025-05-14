package ring

import (
	"unsafe"

	"github.com/matejnesuta/libnf-go/api/errors"
	"github.com/matejnesuta/libnf-go/api/record"
	"github.com/matejnesuta/libnf-go/internal"
)

// Ring represents a ring buffer used for sharing libnf records between processes.
// It is implemented using shared memory. Multiple readers and writers can use the
// same ring buffer concurrently, even across separate processes.
type Ring struct {
	ptr uintptr
}

const (
	// RingTotal is used with Ring.Info() method to retrieve the total number of records
	// properly received since initialization.
	RingTotal int = internal.RING_TOTAL
	// RingLost is used with Ring.Info() method to retrieve the number of records that were lost
	// due to buffer overflows or slow readers.
	RingLost int = internal.RING_LOST
	// RingStuck is used with Ring.Info() method to retrieve the number of times a lock got stuck.
	RingStuck int = internal.RING_STUCK
)

// NewRing initializes and returns a new ring buffer using shared memory.
//
// The filename parameter specifies the name of the shared memory segment (e.g., "libnf-shm").
// The flags control the initialization behavior:
//   - forceInit: reinitializes the shared memory buffer.
//   - forceRelease: removes the shared memory buffer on release.
//   - nonBlockingReading: enables non-blocking reads.
//
// If the process exits without calling Free, the shared memory segment remains
// allocated. It is recommended to use forceInit or forceRelease in at least one
// process (typically the main writer).
//
// Returns a Ring instance or an error if the initialization fails.
func NewRing(filename string, forceInit bool, forceRelease bool, nonBlockingReading bool) (Ring, error) {
	ring := Ring{}
	var flags int
	if forceInit {
		flags |= internal.RING_FORCE_INIT
	}
	if forceRelease {
		flags |= internal.RING_FORCE_RELEASE
	}
	if nonBlockingReading {
		flags |= internal.RING_NO_BLOCK
	}

	status := internal.Ring_init(&ring.ptr, filename, flags)
	if status == internal.ERR_NOMEM {
		return ring, errors.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return ring, errors.ErrOther
	}
	return ring, nil
}

// Info retrieves internal statistics from the ring buffer.
//
// The infoType argument must be one of RingTotal, RingLost, or RingStuck.
// The returned integer holds the corresponding counter value.
// Returns an error if the request fails.
func (r *Ring) Info(infoType int) (int, error) {
	var info uint64
	status := internal.Ring_info(r.ptr, infoType, uintptr(unsafe.Pointer(&info)), int64(8))
	if status == internal.ERR_OTHER {
		return 0, errors.ErrOther
	}
	return status, nil
}

// GetNextRecord reads the next record from the ring buffer into the provided Record.
//
// The caller must ensure the Record is allocated. If no records are available,
// this returns errors.ErrFileEof. If the reader is too slow, some records may be lost.
//
// Use Ring.Info(RingLost) to retrieve the number of lost records.
func (r *Ring) GetNextRecord(rec *record.Record) error {
	if !rec.Allocated() {
		return errors.ErrRecordNotAllocated
	}

	status := internal.Ring_read(r.ptr, rec.GetPtr())
	if status == internal.EOF {
		return errors.ErrFileEof
	} else if status == internal.ERR_OTHER {
		return errors.ErrOther
	}
	return nil
}

// WriteRecord writes the given record into the ring buffer.
//
// The caller must ensure the Record is allocated. Returns an error if memory
// allocation fails or another error occurs.
func (r *Ring) WriteRecord(rec *record.Record) error {
	if !rec.Allocated() {
		return errors.ErrRecordNotAllocated
	}

	status := internal.Ring_write(r.ptr, rec.GetPtr())
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return errors.ErrOther
	}
	return nil
}

// Free releases the ring buffer and its associated resources.
//
// If this is the last instance, and the forceRelease flag was set during initialization,
// the shared memory segment is also removed.
func (r *Ring) Free() error {
	internal.Ring_free(r.ptr)
	return nil
}
