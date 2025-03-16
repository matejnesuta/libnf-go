package examples

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
	heap.SetAggrOptions(fields.FldCalcBps, memheap.AggrKey, memheap.SortDesc, 0, 0)
	heap.SetAggrOptions(fields.FldFirst, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldLast, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldSrcaddr, memheap.AggrAuto, memheap.SortNone, 32, 128)
	heap.SetAggrOptions(fields.FldDstaddr, memheap.AggrAuto, memheap.SortNone, 32, 128)
	heap.SetAggrOptions(fields.FldSrcport, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldDstport, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldProt, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldDoctets, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldDpkts, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldAggrFlows, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldSrcport, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SetAggrOptions(fields.FldDstport, memheap.AggrAuto, memheap.SortNone, 0, 0)

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
		val, _ := rec.GetField(fields.FldBrec1)
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		printBrec(&brec)
	}
}
