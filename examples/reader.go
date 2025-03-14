package examples

import (
	"fmt"
	LnfFld "libnf/api/fields"
	LnfFile "libnf/api/file"
	LnfFilter "libnf/api/filter"
	LnfRec "libnf/api/record"
)

func Reader() {
	var ptr LnfFile.File
	err := ptr.OpenRead("api/testfiles/nfcapd.201705281555", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	var filter LnfFilter.Filter
	filter.Init("src port 80")
	defer filter.Free()
	var num_of_matches uint64 = 0
	fmt.Println("Time first seen|Source IP|Source Port|Destination IP|Destination Port|Pkts|Bytes|Flows")
	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		// if match, _ := filter.Match(rec); !match {
		// 	continue
		// }
		val, _ := rec.GetField(LnfFld.FldBrec1)
		brec, ok := val.(LnfFld.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		// fmt.Print(first.(time.Time).Format("2006-01-02 15:04:05"), " ")
		fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %-15s:%-5d | %-15s:%-5d | %6d | %8d | %4d \n", brec.SrcAddr, brec.SrcPort, brec.DstAddr, brec.DstPort, brec.Pkts, brec.Bytes, brec.Flows)
		num_of_matches++
	}
	fmt.Println("Number of matches captured by '", filter, "' filter: ", num_of_matches)
}
