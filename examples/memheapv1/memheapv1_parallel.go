package examples

import (
	"fmt"
	"libnf/api/fields"
	"libnf/api/file"
	"libnf/api/memheap"
	"libnf/api/record"
	"runtime"
	"sync"
)

func MemHeapV1P() {
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
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	for x := 0; x < 8; x++ {
		wg.Add(1)
		go func() {
			runtime.LockOSThread()         // Lock the goroutine to the OS thread
			defer runtime.UnlockOSThread() // Unlock the goroutine from the OS thread
			defer wg.Done()
			mutex.Lock()
			rec, err := record.NewRecord()
			mutex.Unlock()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer rec.Free()
			for {
				mutex.Lock()
				err = ptr.GetNextRecord(&rec)
				mutex.Unlock()
				if err != nil {
					break
				}
				mutex.Lock()
				i++

				err = heap.WriteRecord(&rec)
				if err != nil {
					fmt.Println(err)
				}
				mutex.Unlock()
			}
			heap.MergeThreads()
		}()
	}
	wg.Wait() // Wait for all workers to finish processing
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
