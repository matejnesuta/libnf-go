package examples

import (
	"fmt"
	"libnf/api/fields"
	"libnf/api/file"
	"libnf/api/memheap"
	"libnf/api/record"
)

func MemHeapV1() {
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

	heap, err := memheap.NewMemHeap()
	if err != nil {
		fmt.Println(err)
	}

	err = heap.SetAggrOptions(fields.SrcAddr, memheap.AggrKey, memheap.SortNone, 24, 64)
	if err != nil {
		panic("uhhhh")
	}
	heap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Last, memheap.AggrMax, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortDesc, 0, 0)
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
	i = 0
	for {
		err = heap.GetNextRecord(&rec)
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
		if err != nil {
			fmt.Println(err)
			break
		}

	}
	fmt.Println("Total records in heap: ", i)
}
