package examples

import (
	"fmt"
	fields "libnf/api/fields"
	file "libnf/api/file"
	record "libnf/api/record"
)

func Reader() {
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

	printHeader()
	for {
		err = ptr.GetNextRecord(&rec)
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
