package examples

import (
	"fmt"
	"libnf/api/file"
	"libnf/api/record"
)

// Writer is a function that demonstrates how to use the libnf package to write data to a file.
func Writer() {
	var output file.File
	err := output.OpenWrite("tmp/writer.tmp", "", false, 0, false)

	if err != nil {
		fmt.Println(err)
	}

	var input file.File
	err = input.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer output.Close()
	defer input.Close()

	rec, err := record.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	// Set the fields of the record
	for {
		err = input.GetNextRecord(&rec)
		if err != nil {
			break
		}

		err = output.WriteRecord(&rec)
		if err != nil {
			fmt.Println(err)
		}
	}
}
