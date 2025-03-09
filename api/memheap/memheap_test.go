package memheap_test

import (
	"libnf/api/memheap"
	"testing"

	LnfErr "libnf/api/errors"
	LnfFld "libnf/api/fields"
	LnfRec "libnf/api/record"

	"github.com/stretchr/testify/assert"
)

func TestNewMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())
}

func TestMemHeapDoubleFree(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, LnfErr.ErrMemHeapNotAllocated, err)
	assert.Equal(t, false, memHeap.Allocated())
}

func TestReadFromEmptyMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	rec, _ := LnfRec.NewRecord()
	err = memHeap.GetNextRecord(&rec)
	assert.Equal(t, LnfErr.ErrMemHeapEnd, err)
}

func TestReadFromUsingUnallocatedRecord(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	rec := LnfRec.Record{}
	err = memHeap.GetNextRecord(&rec)
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}

func TestReadFromFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	rec, _ := LnfRec.NewRecord()
	err = memHeap.GetNextRecord(&rec)
	assert.Equal(t, LnfErr.ErrMemHeapNotAllocated, err)
}

func TestSetAggrOptionsToFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.SetAggrOptions(LnfFld.FldSrcport, memheap.AggrKey, memheap.SortAsc, 0, 0)
	assert.Equal(t, LnfErr.ErrMemHeapNotAllocated, err)
}

func TestWriteToFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	rec, _ := LnfRec.NewRecord()
	err = memHeap.WriteRecord(&rec)
	assert.Equal(t, LnfErr.ErrMemHeapNotAllocated, err)
}

func TestWriteToMemHeap(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	memHeap.SetAggrOptions(LnfFld.FldSrcport, memheap.AggrKey, memheap.SortAsc, 0, 0)

	rec, _ := LnfRec.NewRecord()
	LnfRec.SetField(&rec, LnfFld.FldSrcport, uint16(80))

	err := memHeap.WriteRecord(&rec)
	assert.Equal(t, nil, err)

	rec2, _ := LnfRec.NewRecord()
	err = memHeap.GetNextRecord(&rec2)
	assert.Equal(t, nil, err)

	val, _ := rec2.GetField(LnfFld.FldSrcport)
	srcport := val.(uint16)
	assert.Equal(t, uint16(80), srcport)
}

func TestCleanFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.Clear()
	assert.Equal(t, LnfErr.ErrMemHeapNotAllocated, err)
}

func TestCleanMemHeap(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	memHeap.SetAggrOptions(LnfFld.FldSrcport, memheap.AggrKey, memheap.SortAsc, 0, 0)

	rec, _ := LnfRec.NewRecord()
	LnfRec.SetField(&rec, LnfFld.FldSrcport, uint16(80))

	err := memHeap.WriteRecord(&rec)
	assert.Equal(t, nil, err)

	err = memHeap.Clear()
	assert.Equal(t, nil, err)

	rec2, _ := LnfRec.NewRecord()
	err = memHeap.GetNextRecord(&rec2)
	assert.Equal(t, LnfErr.ErrMemHeapEnd, err)
}
