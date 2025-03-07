package examples

import (
	"fmt"
	LnfFld "libnf/api/fields"
	LnfFile "libnf/api/file"
	LnfFilter "libnf/api/filter"
	LnfRec "libnf/api/record"
	"net"
)

func Reader() {
	var ptr LnfFile.File
	err := ptr.OpenRead("api/tests/testfiles/test-file.tmp", false, false)

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
	fmt.Println("Time first seen | Source IP | Source Port | Destination IP | Destination Port | Out Bytes | ACL ID | ACE ID | XACE ID")
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

		val, _ = rec.GetField(LnfFld.FldIngressAcl)
		acl, ok := val.(LnfFld.Acl)
		if !ok {
			panic("Error: Not an Acl")
		}

		val, _ = rec.GetField(LnfFld.FldMplsLabel)
		mpls, ok := val.(LnfFld.Mpls)
		if !ok {
			panic("Error: Not an Mpls")
		}
		fmt.Println(mpls)

		val, _ = rec.GetField(LnfFld.FldUsername)

		fmt.Println(val)
		// fmt.Println(err)

		val, _ = rec.GetField(LnfFld.FldIpNextHop)
		exporterIp, ok := val.(net.IP)
		if !ok {
			panic("Error: Not an IP address")
		}
		fmt.Println(exporterIp)

		val, _ = rec.GetField(LnfFld.FldSrcAS)
		exporterId, ok := val.(uint32)
		if !ok {
			panic("Error: Not an integer")
		}
		fmt.Println(exporterId)

		val, _ = rec.GetField(LnfFld.FldInSrcMac)
		mac, ok := val.(net.HardwareAddr)
		if !ok {
			panic("Error: Not a MAC address")
		}
		fmt.Println(mac)

		val, _ = rec.GetField(LnfFld.FldCalcBps)
		bps, ok := val.(float64)
		if !ok {
			panic("Error: Not a float64")
		}

		fmt.Println(bps)
		// first, _ := rec.GetField(libnf.FldFirst)
		// srcport, _ := rec.GetField(libnf.FldSrcport)
		// dstport, _ := rec.GetField(libnf.FldDstport)
		// outbytes, _ := rec.GetField(libnf.FldDoctets)

		// srcip, _ := rec.GetField(libnf.FldSrcaddr)
		// dstip, _ := rec.GetField(libnf.FldDstaddr)

		// fmt.Print(first.(time.Time).Format("2006-01-02 15:04:05"), " ")
		fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %s | %d | %s | %d | %d | %d | %d | %d |\n", brec.SrcAddr, brec.SrcPort, brec.DstAddr, brec.DstPort, brec.Bytes, acl.AclId, acl.AceId, acl.XaceId)
		num_of_matches++
	}
	fmt.Println("Number of matches captured by '", filter, "' filter: ", num_of_matches)
}
