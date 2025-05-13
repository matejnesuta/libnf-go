package examples

import (
	"fmt"

	"github.com/matejnesuta/libnf-go/api/fields"
	"github.com/matejnesuta/libnf-go/api/file"
	memheap "github.com/matejnesuta/libnf-go/api/memheapv2"
	"github.com/matejnesuta/libnf-go/api/record"
)

func MemHeapV2() {
	var ptr file.File
	err := ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2(1)

	err = heap.SortAggrOptions(fields.SrcAddr, memheap.AggrKey, memheap.SortNone, 32, 128)
	if err != nil {
		panic("uhhhh")
	}
	heap.SortAggrOptions(fields.SrcPort, memheap.AggrKey, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Last, memheap.AggrMax, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.CalcDuration, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.AggrFlows, memheap.AggrSum, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.CalcBps, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.CalcBpp, memheap.AggrAuto, memheap.SortAsc, 0, 0)
	heap.SortAggrOptions(fields.CalcPps, memheap.AggrAuto, memheap.SortNone, 0, 0)

	rec, _ := record.NewRecord()
	defer rec.Free()

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
		val, _ = rec.GetField(fields.CalcBpp)
		bpp, ok := val.(float64)
		if !ok {
			panic("Error: Not a CalcBpp")
		}
		val, _ = rec.GetField(fields.CalcBps)
		bps, ok := val.(float64)
		if !ok {
			panic("Error: Not a CalcBps")
		}
		val, _ = rec.GetField(fields.CalcPps)
		pps, ok := val.(float64)
		if !ok {
			panic("Error: Not a CalcPps")
		}
		i++
		fmt.Println(brec.First.Format("2006-01-02 15:04:05"), brec.Last.Sub(brec.First).Seconds(), brec.SrcAddr, brec.Bytes, brec.Pkts, brec.Flows, bpp, bps, pps)
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			break
		}
	}
	fmt.Println("Total records in heap: ", i)
}
