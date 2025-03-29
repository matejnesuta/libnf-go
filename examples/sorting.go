package examples

//this example file should currently match the one of the original libnf: https://github.com/netx-as/libnf/blob/master/examples/lnf_ex03_aggreg.c

import (
	"fmt"
	"libnf/api/fields"
	"libnf/api/file"
	"libnf/api/memheap"
	"libnf/api/record"
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

	heap, err := memheap.NewMemHeap()
	if err != nil {
		fmt.Println(err)
	}
	defer heap.Free()

	// heap.Clear()
	// heap.EnableNfdumpCompat()
	// heap.SetHashBuckets(2000000000)
	// heap.SetListMode()
	err = heap.SetAggrOptions(fields.SrcAddr, memheap.AggrKey, memheap.SortDesc, 24, 64)
	if err != nil {
		panic("uhhhh")
	}
	heap.SetAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.DstAS, memheap.AggrKey, memheap.SortNone, 0, 0)

	heap.SetAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Last, memheap.AggrMax, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.TcpFlags, memheap.AggrOr, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
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
			panic(err)
		}
		i++
		if i == 1000000 {
			break
		}
	}
	fmt.Println("Total records in file: ", i)
	// printHeader()
	cursor, _ := heap.FirstRecordPosition()
	i = 0
	for {
		if i == 400000 {
			fmt.Println("break")
		}
		// err = heap.GetNextRecord(&rec)
		err = heap.GetRecordWithCursor(&cursor, &rec)
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
		err = heap.NextRecordPosition(&cursor)
		if err != nil {
			fmt.Println(err)
			break
		}

	}
	fmt.Println("Total records in heap: ", i)
}
