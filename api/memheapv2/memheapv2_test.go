package memheapv2_test

import (
	"libnf/api/errors"
	"libnf/api/fields"
	memheap "libnf/api/memheapv2"
	"libnf/api/record"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sortByUint8(t *testing.T, ports []uint16, protocols []uint8, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Prot, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [3]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(111),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("1.1.1.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.3").To4(),
		DstAddr: net.ParseIP("1.1.1.4").To4(),
		Prot:    uint8(5),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(90),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.5").To4(),
		DstAddr: net.ParseIP("1.1.1.6").To4(),
		Prot:    uint8(4),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

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
		assert.Equal(t, protocols[i], brec.Prot)
		i++
	}
}

func sortByUint16(t *testing.T, bytes []uint64, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.Doctets, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [3]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(111),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("1.1.1.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.3").To4(),
		DstAddr: net.ParseIP("1.1.1.4").To4(),
		Prot:    uint8(5),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(90),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.5").To4(),
		DstAddr: net.ParseIP("1.1.1.6").To4(),
		Prot:    uint8(4),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

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
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
}

func TestSortByUint8Asc(t *testing.T) {
	ports := [3]uint16{90, 80, 111}
	protocols := [3]uint8{4, 5, 6}
	sortByUint8(t, ports[:], protocols[:], memheap.SortAsc)
}

func TestSortByUint8Desc(t *testing.T) {
	ports := [3]uint16{111, 80, 90}
	protocols := [3]uint8{6, 5, 4}
	sortByUint8(t, ports[:], protocols[:], memheap.SortDesc)
}

func TestSortByUint16Asc(t *testing.T) {
	bytes := [3]uint64{80, 80, 20}
	ports := [3]uint16{80, 90, 111}
	sortByUint16(t, bytes[:], ports[:], memheap.SortAsc)
}

func TestSortByUint16Desc(t *testing.T) {
	bytes := [3]uint64{20, 80, 80}
	ports := [3]uint16{111, 90, 80}
	sortByUint16(t, bytes[:], ports[:], memheap.SortDesc)
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
