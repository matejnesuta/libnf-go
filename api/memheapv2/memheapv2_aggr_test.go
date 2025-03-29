package memheapv2_test

import (
	"libnf/api/errors"
	"libnf/api/fields"
	"libnf/api/file"
	memheap "libnf/api/memheapv2"
	"libnf/api/record"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAggrMinMaxField(t *testing.T) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Last, memheap.AggrMax, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	first := []time.Time{time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local),
		time.Date(2017, time.May, 28, 15, 55, 20, 0, time.Local),
		time.Date(2017, time.May, 28, 15, 55, 40, 0, time.Local)}
	last := []time.Time{time.Date(2017, time.May, 28, 15, 56, 0, 0, time.Local),
		time.Date(2017, time.May, 28, 15, 56, 20, 0, time.Local),
		time.Date(2017, time.May, 28, 15, 56, 40, 0, time.Local)}
	doctets := []uint64{20, 40, 60}

	for i := 0; i < 3; i++ {
		record.SetField(&rec, fields.First, first[i])
		record.SetField(&rec, fields.Last, last[i])
		record.SetField(&rec, fields.Doctets, doctets[i])
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)

	err = heap.GetRecord(&cursor, &rec)
	assert.Nil(t, err)
	cursor, err = heap.NextRecordPosition(cursor)
	assert.Equal(t, err, errors.ErrMemHeapEnd)
	val, _ := rec.GetField(fields.First)
	assert.Equal(t, time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local), val)
	val, _ = rec.GetField(fields.Last)
	assert.Equal(t, time.Date(2017, time.May, 28, 15, 56, 40, 0, time.Local), val)
	val, _ = rec.GetField(fields.Doctets)
	assert.Equal(t, uint64(120), val)
}

func TestAggrPerPairField(t *testing.T) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SetNfdumpComp(false)

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
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	ports := [3]uint16{53, 80, 1222}
	pkts := [3]uint64{2, 3, 3}
	bytes := [3]uint64{40, 80, 80}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, pkts[i], brec.Pkts)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
}

func TestAggrPerPairFieldWithNfdumpComp(t *testing.T) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SetNfdumpComp(true)

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
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	ports := [3]uint16{53, 80, 1222}
	pkts := [3]uint64{1, 3, 3}
	bytes := [3]uint64{20, 80, 80}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, pkts[i], brec.Pkts)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
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

	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()

	err = heap.SortAggrOptions(fields.PairAddr, memheap.AggrKey, memheap.SortNone, 24, 64)
	assert.Equal(t, nil, err)
	err = heap.SortAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)
	err = heap.SortAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)
	err = heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortDesc, 0, 0)
	assert.Equal(t, nil, err)
	err = heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Equal(t, nil, err)

	for {
		err = file.GetNextRecord(&rec)
		if err != nil {
			break
		}
		records++
		err = heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, 2035, records)

	var doctets = [10]uint64{38232, 12773, 12721, 12514, 9959, 3988, 3988, 525, 400, 328}
	var packets = [10]uint64{64, 22, 21, 17, 94, 29, 29, 9, 8, 1}

	cursor, err := heap.FirstRecordPosition()

	assert.Nil(t, err)
	records = 0
	for {
		err = heap.GetRecord(&cursor, &rec)
		if err != nil {
			break
		}
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			break
		}
		if records < 10 {
			val, _ := rec.GetField(fields.Doctets)
			assert.Equal(t, doctets[records], val)
			val, _ = rec.GetField(fields.Dpkts)
			assert.Equal(t, packets[records], val)
		}

		records++
	}

	assert.Equal(t, 1982, records+1)
}
