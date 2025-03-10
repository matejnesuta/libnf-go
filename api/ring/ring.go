package ring

import (
	LnfErr "libnf/api/errors"
	LnfRec "libnf/api/record"
	"libnf/internal"
	"unsafe"
)

type Ring struct {
	ptr uintptr
}

const (
	RingTotal int = internal.RING_TOTAL
	RingLost  int = internal.RING_LOST
	RingStuck int = internal.RING_STUCK
)

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
		return ring, LnfErr.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return ring, LnfErr.ErrOther
	}
	return ring, nil
}

func (r *Ring) Info(infoType int) (int, error) {
	var info uint64
	status := internal.Ring_info(r.ptr, infoType, uintptr(unsafe.Pointer(&info)), int64(8))
	if status == internal.ERR_OTHER {
		return 0, LnfErr.ErrOther
	}
	return status, nil
}

func (r *Ring) GetNextRecord(rec *LnfRec.Record) error {
	if !rec.Allocated() {
		return LnfErr.ErrRecordNotAllocated
	}

	status := internal.Ring_read(r.ptr, rec.GetPtr())
	if status == internal.EOF {
		return LnfErr.ErrFileEof
	} else if status == internal.ERR_OTHER {
		return LnfErr.ErrOther
	}
	return nil
}

func (r *Ring) WriteRecord(rec *LnfRec.Record) error {
	if !rec.Allocated() {
		return LnfErr.ErrRecordNotAllocated
	}

	status := internal.Ring_write(r.ptr, rec.GetPtr())
	if status == internal.ERR_NOMEM {
		return LnfErr.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return LnfErr.ErrOther
	}
	return nil
}

func (r *Ring) Free() error {
	internal.Ring_free(r.ptr)
	return nil
}
