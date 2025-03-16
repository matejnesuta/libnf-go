package examples

// nfdump -r api/testfiles/profiling.tmp -n 10  -O bps

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

	heap.Clear()
	heap.EnableNfdumpCompat()
	heap.SetAggrOptions(fields.CalcBps, memheap.AggrAuto, memheap.SortDesc, 0, 0)
	heap.SetAggrOptions(fields.First, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Last, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.SrcAddr, memheap.AggrAuto, memheap.SortNone, 32, 128)
	heap.SetAggrOptions(fields.DstAddr, memheap.AggrAuto, memheap.SortNone, 32, 128)
	heap.SetAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.DstPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Prot, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Doctets, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.Dpkts, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.AggrFlows, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.DstPort, memheap.AggrAuto, memheap.SortNone, 0, 0)

	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		heap.WriteRecord(&rec)
	}

	printHeader()

	for i := 0; i < 10; i++ {
		err = heap.GetNextRecord(&rec)
		if err != nil {
			break
		}
		val, _ := rec.GetField(fields.Brec1)
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		printBrec(&brec)
	}
}
