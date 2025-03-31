package examples

import (
	"fmt"
	"libnf/api/fields"
	"libnf/api/file"
	memheap "libnf/api/memheapv2"
	"libnf/api/record"
	"sync"
)

func Psorting() {
	var ptr file.File
	err := ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	heap := memheap.NewMemHeapV2()
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

	var i uint64 = 0
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	for x := 0; x < 16; x++ {
		wg.Add(1)
		go func() {
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

				err = heap.WriteRecord(&rec)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
	wg.Wait() // Wait for all workers to finish processing
	fmt.Println("Total records in file: ", i)
	cursor, _ := heap.FirstRecordPosition()
	i = 0
	rec, _ := record.NewRecord()
	defer rec.Free()
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
