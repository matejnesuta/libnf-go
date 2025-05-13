package examples

import (
	"fmt"
	"libnf-go/api/fields"
	"libnf-go/api/file"
	memheap "libnf-go/api/memheapv2"
	"libnf-go/api/record"
)

func Sorting() {
	var ptr file.File
	err := ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := record.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	heap := memheap.NewMemHeapV2(1)
	if err != nil {
		fmt.Println(err)
	}

	err = heap.SortAggrOptions(fields.SrcAddr, memheap.AggrKey, memheap.SortNone, 24, 64)
	if err != nil {
		panic("uhhhh")
	}
	heap.SortAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Last, memheap.AggrMax, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortDesc, 0, 0)
	// heap.SetAggrOptions(fields.CalcBps, memheap.AggrAuto, memheap.SortDesc, 0, 0)

	var i uint64 = 0
	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		err = heap.WriteRecord(&rec)
		if err != nil {
			fmt.Println(err)
		}
		i++
	}
	fmt.Println("Total records in file: ", i)
	cursor, _ := heap.FirstRecordPosition()
	i = 0
	for {
		err = heap.GetRecord(&cursor, &rec)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error getting record")
			break
		}
		val, err := rec.GetField(fields.Brec1)
		if err != nil {
			panic(err)
		}
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}
		i++
		fmt.Println(brec.SrcAddr, brec.SrcPort, brec.DstAddr, brec.DstPort, brec.First.UnixMilli(), brec.Bytes, brec.Pkts)
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			fmt.Println(err)
			break
		}

	}
	fmt.Println("Total records in heap: ", i)
}
