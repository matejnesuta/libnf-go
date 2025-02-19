package main

import (
	"fmt"
	libnf "libnf/api"
)

func main() {
	var ptr libnf.File
	err := ptr.OpenRead("api/tests/testfiles/nfcapd.201705281555", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := libnf.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	var filter libnf.Filter
	filter.Init("src port 80")
	defer filter.Free()
	var num_of_matches uint64 = 0
	fmt.Println("Time first seen | Source IP | Source Port | Destination IP | Destination Port | Out Bytes")
	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		// if match, _ := filter.Match(rec); !match {
		// 	continue
		// }
		val, _ := rec.GetField(libnf.FldBrec1)
		brec, ok := val.(libnf.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		// first, _ := rec.GetField(libnf.FldFirst)
		// srcport, _ := rec.GetField(libnf.FldSrcport)
		// dstport, _ := rec.GetField(libnf.FldDstport)
		// outbytes, _ := rec.GetField(libnf.FldDoctets)

		// srcip, _ := rec.GetField(libnf.FldSrcaddr)
		// dstip, _ := rec.GetField(libnf.FldDstaddr)

		// fmt.Print(first.(time.Time).Format("2006-01-02 15:04:05"), " ")
		fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %s | %d | %s | %d | %d\n", brec.SrcAddr, brec.SrcPort, brec.DstAddr, brec.DstPort, brec.Bytes)
		num_of_matches++
	}
	fmt.Println("Number of matches captured by '", filter, "' filter: ", num_of_matches)
}
