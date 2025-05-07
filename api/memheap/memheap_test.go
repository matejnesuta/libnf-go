package memheap_test

import (
	"libnf/api/memheap"
	"net"
	"testing"
	"time"

	"libnf/api/errors"
	"libnf/api/fields"
	"libnf/api/file"
	"libnf/api/record"

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
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
	assert.Equal(t, false, memHeap.Allocated())
}

func TestReadFromEmptyMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	rec, _ := record.NewRecord()
	err = memHeap.GetNextRecord(&rec)
	assert.Equal(t, errors.ErrMemHeapEnd, err)
}

func TestReadFromUsingUnallocatedRecord(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	rec := record.Record{}
	err = memHeap.GetNextRecord(&rec)
	assert.Equal(t, errors.ErrRecordNotAllocated, err)
}

func TestReadFromFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	rec, _ := record.NewRecord()
	err = memHeap.GetNextRecord(&rec)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestSetAggrOptionsToFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestWriteToFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	rec, _ := record.NewRecord()
	err = memHeap.WriteRecord(&rec)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestWriteUnallocatedRecordToMemHeap(t *testing.T) {
	memheap, _ := memheap.NewMemHeap()
	defer memheap.Free()

	rec := record.Record{}
	err := memheap.WriteRecord(&rec)
	assert.Equal(t, errors.ErrRecordNotAllocated, err)
}

func TestWriteToMemHeap(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	memHeap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	defer memHeap.Free()

	rec, _ := record.NewRecord()
	defer rec.Free()
	record.SetField(&rec, fields.SrcPort, uint16(80))

	err := memHeap.WriteRecord(&rec)
	assert.Equal(t, nil, err)

	rec2, _ := record.NewRecord()
	defer rec2.Free()
	err = memHeap.GetNextRecord(&rec2)
	assert.Equal(t, nil, err)

	val, _ := rec2.GetField(fields.SrcPort)
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
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestCleanMemHeap(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	memHeap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)

	rec, _ := record.NewRecord()
	record.SetField(&rec, fields.SrcPort, uint16(80))

	err := memHeap.WriteRecord(&rec)
	assert.Equal(t, nil, err)

	err = memHeap.Clear()
	assert.Equal(t, nil, err)

	rec2, _ := record.NewRecord()
	err = memHeap.GetNextRecord(&rec2)
	assert.Equal(t, errors.ErrMemHeapEnd, err)
}

func TestSetListModeOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.SetListMode()
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestSetHashBucketsOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.SetHashBuckets(10)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestEnableNfdumpCompatOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.EnableNfdumpCompat()
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestSetFastAggrOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.SetFastAggr(memheap.FastAggrNone)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestMergeThreadsOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	err = memHeap.MergeThreads()
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestFirstRecordPositionOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	_, err = memHeap.FirstRecordPosition()
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestNextRecordPositionOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	cursor, _ := memHeap.FirstRecordPosition()
	err = memHeap.NextRecordPosition(&cursor)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestGetRecordWithCursorOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	cursor, _ := memHeap.FirstRecordPosition()
	rec, _ := record.NewRecord()
	err = memHeap.GetRecordWithCursor(&cursor, &rec)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestFirstRecordPositionOnEmptyMemHeap(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	defer memHeap.Free()

	_, err := memHeap.FirstRecordPosition()
	assert.Equal(t, errors.ErrMemHeapEnd, err)
}

func TestNextRecordPositionOnEmptyMemHeap(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	defer memHeap.Free()

	var cursor memheap.MemHeapCursor
	err := memHeap.NextRecordPosition(&cursor)
	assert.Equal(t, errors.ErrMemHeapEnd, err)
}

func TestGetRecordWithCursorOnUnallocatedRecord(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	memHeap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)

	cursor, _ := memHeap.FirstRecordPosition()
	rec := record.Record{}
	err := memHeap.GetRecordWithCursor(&cursor, &rec)
	assert.Equal(t, errors.ErrRecordNotAllocated, err)
}

func TestGetRecordWithKeyOnFreedMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())

	rec, _ := record.NewRecord()
	_, err = memHeap.GetRecordWithKey(&rec)
	assert.Equal(t, errors.ErrMemHeapNotAllocated, err)
}

func TestGetRecordWithKeyOnUnallocatedRecord(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	memHeap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)

	rec := record.Record{}
	_, err := memHeap.GetRecordWithKey(&rec)
	assert.Equal(t, errors.ErrRecordNotAllocated, err)
}

func TestStatistics(t *testing.T) {
	var file = file.File{}
	err := file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	defer file.Close()
	assert.Equal(t, nil, err)

	rec, err := record.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	var records int = 0

	memHeap, err := memheap.NewMemHeap()
	defer memHeap.Free()
	assert.Equal(t, nil, err)

	err = memHeap.SetAggrOptions(fields.PairAddr, memheap.AggrKey, memheap.SortNone, 24, 64)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortDesc, 0, 0)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)

	for {
		err = file.GetNextRecord(&rec)
		if err != nil {
			break
		}
		records++
		err = memHeap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, 2035, records)

	var doctets [2]uint64
	memHeap.GetNextRecord(&rec)
	val, err := rec.GetField(fields.Doctets)
	assert.Equal(t, nil, err)
	doctets[0] = val.(uint64)
	memHeap.GetNextRecord(&rec)
	val, err = rec.GetField(fields.Doctets)
	assert.Equal(t, nil, err)
	doctets[1] = val.(uint64)

	assert.Equal(t, uint64(38232), doctets[0])
	assert.Equal(t, uint64(12773), doctets[1])

	records = 2
	for {
		err = memHeap.GetNextRecord(&rec)
		if err != nil {
			break
		}
		records++
	}

	assert.Equal(t, 1982, records)
}

func TestStatisticsListMode(t *testing.T) {
	var file = file.File{}
	err := file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	defer file.Close()
	assert.Equal(t, nil, err)

	rec, err := record.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	var records int = 0

	memHeap, err := memheap.NewMemHeap()
	defer memHeap.Free()
	assert.Equal(t, nil, err)

	err = memHeap.SetListMode()
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.PairAddr, memheap.AggrKey, memheap.SortNone, 24, 64)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortDesc, 0, 0)
	assert.Equal(t, nil, err)
	err = memHeap.SetAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)

	for {
		err = file.GetNextRecord(&rec)
		if err != nil {
			break
		}
		records++
		err = memHeap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, 2035, records)

	var doctets [2]uint64
	memHeap.GetNextRecord(&rec)
	val, err := rec.GetField(fields.Doctets)
	assert.Equal(t, nil, err)
	doctets[0] = val.(uint64)
	memHeap.GetNextRecord(&rec)
	val, err = rec.GetField(fields.Doctets)
	assert.Equal(t, nil, err)
	doctets[1] = val.(uint64)

	assert.Equal(t, uint64(12227), doctets[0])
	assert.Equal(t, uint64(12227), doctets[1])

	records = 2
	for {
		err = memHeap.GetNextRecord(&rec)
		if err != nil {
			break
		}
		records++
	}

	assert.Equal(t, 2035, records)
}

func TestCursor(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	defer memHeap.Free()
	memHeap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)

	rec, _ := record.NewRecord()
	defer rec.Free()
	record.SetField(&rec, fields.SrcPort, uint16(80))

	err := memHeap.WriteRecord(&rec)
	assert.Equal(t, nil, err)

	record.SetField(&rec, fields.SrcPort, uint16(90))

	err = memHeap.WriteRecord(&rec)
	assert.Equal(t, nil, err)

	cursor, err := memHeap.FirstRecordPosition()
	assert.Equal(t, nil, err)

	err = memHeap.GetRecordWithCursor(&cursor, &rec)
	assert.Equal(t, nil, err)

	val, _ := rec.GetField(fields.SrcPort)
	srcport := val.(uint16)
	assert.Equal(t, uint16(80), srcport)

	err = memHeap.NextRecordPosition(&cursor)
	assert.Equal(t, nil, err)

	err = memHeap.GetRecordWithCursor(&cursor, &rec)
	assert.Equal(t, nil, err)

	val, _ = rec.GetField(fields.SrcPort)
	srcport = val.(uint16)
	assert.Equal(t, uint16(90), srcport)
}

func TestNfdumpCompMode(t *testing.T) {
	memHeap, _ := memheap.NewMemHeap()
	defer memHeap.Free()
	memHeap.SetAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	memHeap.SetAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	memHeap.SetAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	err := memHeap.EnableNfdumpCompat()
	assert.Equal(t, nil, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [2]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(53),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("2.2.2.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("3.3.3.3").To4(),
		DstAddr: net.ParseIP("4.4.4.4").To4(),
		Prot:    uint8(6),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := memHeap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	ports := [3]uint16{53, 80, 1222}
	pkts := [3]uint64{1, 3, 3}
	bytes := [3]uint64{20, 80, 80}

	i := 0
	for {
		err := memHeap.GetNextRecord(&rec)
		if err != nil {
			break
		}
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, pkts[i], brec.Pkts)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
	keyRec, _ := record.NewRecord()
	record.SetField(&keyRec, fields.SrcPort, uint16(53))
	cursor, err := memHeap.GetRecordWithKey(&keyRec)
	assert.Equal(t, nil, err)

	memHeap.GetRecordWithCursor(&cursor, &rec)
	val, _ := rec.GetField(fields.SrcPort)
	srcport := val.(uint16)
	assert.Equal(t, uint16(53), srcport)

	val, _ = rec.GetField(fields.Dpkts)
	dpkts := val.(uint64)
	assert.Equal(t, uint64(1), dpkts)
}
