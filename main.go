package main

import (
	"fmt"
	libnf "libnf/api"
	"time"
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
			// fmt.Println(err)
			break
		}
		if match, _ := filter.Match(rec); !match {
			continue
		}
		first, _ := rec.GetField(libnf.FldFirst)
		srcport, _ := rec.GetField(libnf.FldSrcport)
		dstport, _ := rec.GetField(libnf.FldDstport)
		outbytes, _ := rec.GetField(libnf.FldDoctets)

		srcip, _ := rec.GetField(libnf.FldSrcaddr)
		dstip, _ := rec.GetField(libnf.FldDstaddr)
		// srcmac, _ := rec.GetField(libnf.FldInSrcMac)
		// dstmac, _ := rec.GetField(libnf.FldInDstMac)
		fmt.Print(first.(time.Time).Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %s | %d | %s | %d | %d\n", srcip, srcport, dstip, dstport, outbytes)
		num_of_matches++
		// fmt.Println(srcmac.(net.HardwareAddr).String(), dstmac.(net.HardwareAddr).String())
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("srcport: %d, dstport: %d\n", srcport, dstport)
	}
	fmt.Println("Number of matches captured by '", filter, "' filter: ", num_of_matches)
}
