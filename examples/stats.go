package examples

// nfdump -r api/testfiles/profiling.tmp -s srcip/bps -s dstip/bps  ' port 443'

import (
	"fmt"

	"github.com/matejnesuta/libnf-go/api/fields"
	"github.com/matejnesuta/libnf-go/api/file"
	"github.com/matejnesuta/libnf-go/api/filter"
	memheap "github.com/matejnesuta/libnf-go/api/memheapv2"
	"github.com/matejnesuta/libnf-go/api/record"
)

func Stats() {
	var ptr file.File
	err := ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		panic(err)
		// fmt.Println(err)
	}

	rec, err := record.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	filter := filter.Filter{}
	err = filter.Init("port 443")
	if err != nil {
		fmt.Println(err)
	}

	heap := memheap.NewMemHeapV2(1)
	// heap, err := memheap.NewMemHeap()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// heap.Clear()
	heap.SetNfdumpComp(true)
	heap.SortAggrOptions(fields.CalcBps, memheap.AggrAuto, memheap.SortDesc, 0, 0)
	heap.SortAggrOptions(fields.CalcBpp, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.CalcPps, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.First, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Last, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.SrcAddr, memheap.AggrKey, memheap.SortNone, 32, 128)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Dpkts, memheap.AggrAuto, memheap.SortNone, 0, 0)

	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}

		if match, _ := filter.Match(rec); !match {
			continue
		}

		heap.WriteRecord(&rec)
	}
	fmt.Println("Top 10 Src IP Addr ordered by bps:")
	fmt.Println("")
	fmt.Print("First\t\tDuration\tSrcAddr\t\tPackets\tBytes\t\tPps\t\tBps\t\tBpp\n")

	cursor, err := heap.FirstRecordPosition()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		err = heap.GetRecord(&cursor, &rec)
		if err != nil {
			break
		}
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			break
		}
		val, _ := rec.GetField(fields.Brec1)
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		val, _ = rec.GetField(fields.CalcBps)
		bps, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		val, _ = rec.GetField(fields.CalcPps)
		pps, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		val, _ = rec.GetField(fields.CalcBpp)
		bpp, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %.3f | %-15s| %4d | %4d | %4f | %4f | %4f \n", brec.Last.Sub(brec.First).Seconds(), brec.SrcAddr, brec.Pkts, brec.Bytes, pps, bps, bpp)
	}

	ptr.Close()
	err = ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	heap.Clear()
	heap.SetNfdumpComp(true)
	heap.SortAggrOptions(fields.CalcBps, memheap.AggrAuto, memheap.SortDesc, 0, 0)
	heap.SortAggrOptions(fields.CalcBpp, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.CalcPps, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.First, memheap.AggrMin, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Last, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.DstAddr, memheap.AggrKey, memheap.SortNone, 32, 128)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrAuto, memheap.SortNone, 0, 0)
	heap.SortAggrOptions(fields.Dpkts, memheap.AggrAuto, memheap.SortNone, 0, 0)

	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		if match, _ := filter.Match(rec); !match {
			continue
		}
		heap.WriteRecord(&rec)
	}
	fmt.Println("")
	fmt.Println("Top 10 Dst IP Addr ordered by bps:")
	fmt.Println("")
	fmt.Print("First\t\tDuration\tSrcAddr\t\tPackets\tBytes\t\tPps\t\tBps\t\tBpp\n")

	cursor, err = heap.FirstRecordPosition()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 10; i++ {
		err = heap.GetRecord(&cursor, &rec)
		if err != nil {
			break
		}
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			break
		}

		val, _ := rec.GetField(fields.Brec1)
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		val, _ = rec.GetField(fields.CalcBps)
		bps, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		val, _ = rec.GetField(fields.CalcPps)
		pps, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		val, _ = rec.GetField(fields.CalcBpp)
		bpp, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %.3f | %-15s| %4d | %4d | %4f | %4f | %4f \n", brec.Last.Sub(brec.First).Seconds(), brec.DstAddr, brec.Pkts, brec.Bytes, pps, bps, bpp)
	}
}
