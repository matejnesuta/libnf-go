package examples

import (
	"fmt"
	"libnf/api/fields"
	memheap "libnf/api/memheapv2"
	"libnf/api/record"
	"net"
	"time"
)

func MemHeapV2() {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2(1)
	heap.SortAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SetNfdumpComp(false)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [2]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(80),
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
		if err != nil {
			panic(err)
		}
	}

	cursor, err := heap.FirstRecordPosition()
	if err != nil {
		panic(err)
	}
	for {

		err := heap.GetRecord(&cursor, &rec)
		if err != nil {
			break
		}
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		fmt.Println(brec.SrcPort, brec.DstPort, brec.Bytes, brec.Pkts)
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			break
		}
	}
}
