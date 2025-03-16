package examples

import (
	"fmt"
	LnfFld "libnf/api/fields"
	LnfFile "libnf/api/file"
	LnfRec "libnf/api/record"
)

func Reader() {
	var ptr LnfFile.File
	err := ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
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
		val, _ := rec.GetField(LnfFld.FldBrec1)
		brec, ok := val.(LnfFld.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}
		printBrec(&brec)
	}
}
