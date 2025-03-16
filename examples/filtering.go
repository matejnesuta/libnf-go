package examples

import (
	"fmt"
	"libnf/api/fields"
	"libnf/api/file"
	"libnf/api/filter"
	"libnf/api/record"
)

func Filtering() {
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

	var filter filter.Filter
	filter.Init("dst port 53 and proto tcp")
	defer filter.Free()
	var num_of_matches uint64 = 0
	printHeader()
	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		if match, _ := filter.Match(rec); !match {
			continue
		}
		val, _ := rec.GetField(fields.FldBrec1)
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		printBrec(&brec)
		num_of_matches++
	}
	fmt.Println("Number of matches captured by '", filter, "' filter: ", num_of_matches)
}
