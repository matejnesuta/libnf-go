package main

import (
	"fmt"
	libnf "libnf/api"
	"net"
	"time"
)

func main() {
	var ptr libnf.File
	err := ptr.OpenRead("api/tests/testfiles/nfcapd.201705281555", false, false)
	ident, err := ptr.GetIdent()

	fmt.Println(ident)
	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := libnf.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()
	fmt.Println("Time first seen | Source IP | Source Port | Destination IP | Destination Port | Out Bytes")
	// for {
	err = ptr.GetNextRecord(&rec)
	if err != nil {
		fmt.Println(err)
		// break
	}
	first, _ := rec.GetField(libnf.FldFirst)
	srcport, _ := rec.GetField(libnf.FldSrcport)
	dstport, _ := rec.GetField(libnf.FldDstport)
	outbytes, _ := rec.GetField(libnf.FldOutBytes)

	srcip, _ := rec.GetField(libnf.FldSrcaddr)
	dstip, _ := rec.GetField(libnf.FldDstaddr)
	srcmac, _ := rec.GetField(libnf.FldInSrcMac)
	dstmac, _ := rec.GetField(libnf.FldInDstMac)
	fmt.Print(first.(time.Time).Format("2006-01-02 15:04:05"), " ")
	fmt.Printf("| %s | %d | %s | %d | %d\n", srcip, srcport, dstip, dstport, outbytes)
	fmt.Println(srcmac.(net.HardwareAddr).String(), dstmac.(net.HardwareAddr).String())
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("srcport: %d, dstport: %d\n", srcport, dstport)
	// }
}
