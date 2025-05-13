package examples

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"github.com/matejnesuta/libnf-go/api/fields"
	"github.com/matejnesuta/libnf-go/api/file"
	memheap "github.com/matejnesuta/libnf-go/api/memheapv2"
	"github.com/matejnesuta/libnf-go/api/record"
)

func MemHeapV2P() {

	rec, err := record.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2(256)
	if err != nil {
		fmt.Println(err)
	}

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

	var i uint64 = 0
	var wg sync.WaitGroup
	recordMux := &sync.Mutex{}
	incrementMux := &sync.Mutex{}

	for x := 1; x < 6; x++ {
		wg.Add(1)
		go func() {
			runtime.LockOSThread()         // Lock the goroutine to the OS thread
			defer runtime.UnlockOSThread() // Unlock the goroutine from the OS thread
			defer wg.Done()
			var ptr file.File
			err := ptr.OpenRead("api/testfiles/comparison/"+strconv.Itoa(x)+".tmp", false, false)

			if err != nil {
				fmt.Println(err)
			}
			defer ptr.Close()
			recordMux.Lock()
			rec, err := record.NewRecord()
			recordMux.Unlock()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer rec.Free()
			for {
				err = ptr.GetNextRecord(&rec)
				if err != nil {
					break
				}
				incrementMux.Lock()
				i++
				incrementMux.Unlock()
				err = heap.WriteRecord(&rec)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
	wg.Wait() // Wait for all workers to finish processing
	fmt.Println("Total records in file: ", i)
	i = 0
	cursor, err := heap.FirstRecordPosition()
	if err != nil {
		panic(err)
	}
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
		if err != nil {
			fmt.Println(err)
			break
		}
		cursor, err = heap.NextRecordPosition(cursor)
		if err != nil {
			break
		}

	}
	fmt.Println("Total records in heap: ", i)
}
