package examples

import (
	"fmt"
	"libnf/api/fields"
)

func printBrec(brec *fields.BasicRecord1) {
	fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
	fmt.Printf("| %.3f | %3d | %-15s:%-5d | %-15s:%-5d | %6d | %8d | %4d \n", brec.Last.Sub(brec.First).Seconds(), brec.Prot, brec.SrcAddr, brec.SrcPort, brec.DstAddr, brec.DstPort, brec.Pkts, brec.Bytes, brec.Flows)
}

func printHeader() {
	fmt.Println("Time first seen|Duration|Protocol|Source IP|Source Port|Destination IP|Destination Port|Pkts|Bytes|Flows")
}
