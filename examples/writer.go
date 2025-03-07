package examples

import (
	"fmt"
	LnfFld "libnf/api/fields"
	LnfFile "libnf/api/file"
	LnfRec "libnf/api/record"
	"net"
	"time"
)

// Writer is a function that demonstrates how to use the libnf package to write data to a file.
func Writer() {
	var ptr LnfFile.File
	err := ptr.OpenWrite("tmp/writer.tmp", "", false, 0, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	// Set the fields of the record

	t := time.Now()
	brec := LnfFld.BasicRecord1{
		First:   t,
		Last:    t,
		SrcAddr: net.IPv4(192, 168, 1, 1),
		DstAddr: net.IPv4(192, 168, 1, 2),
		Prot:    6,
		SrcPort: 80,
		DstPort: 80,
		Bytes:   100,
		Pkts:    1,
		Flows:   10,
	}

	acl := LnfFld.Acl{
		AclId:  1,
		AceId:  2,
		XaceId: 3,
	}

	mpls := LnfFld.Mpls{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	username := "leopold kasing"

	mac := net.HardwareAddr{0x00, 0x08, 0x00, 0x0c, 0x02, 0x01}

	exporterIp := net.IPv4(10, 10, 0, 1)
	var exporterId uint32 = 10

	LnfRec.SetField(&rec, LnfFld.FldSrcAS, exporterId)
	LnfRec.SetField(&rec, LnfFld.FldBrec1, brec)
	LnfRec.SetField(&rec, LnfFld.FldIngressAcl, acl)
	LnfRec.SetField(&rec, LnfFld.FldMplsLabel, mpls)
	LnfRec.SetField(&rec, LnfFld.FldUsername, username)
	LnfRec.SetField(&rec, LnfFld.FldInSrcMac, mac)

	LnfRec.SetField(&rec, LnfFld.FldIpNextHop, exporterIp)

	err = ptr.WriteRecord(&rec)

	if err != nil {
		fmt.Println(err)
	}

	srcIpBytes := [16]byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	dstIpBytes := [16]byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}
	exporterIpBytes := [16]byte{0x0a, 0x0a, 0x00, 0x01}

	srcIp := net.IP(srcIpBytes[:])
	dstIp := net.IP(dstIpBytes[:])

	brec.SrcAddr = srcIp
	brec.DstAddr = dstIp
	LnfRec.SetField(&rec, LnfFld.FldExporterIp, net.IP(exporterIpBytes[:]))
	LnfRec.SetField(&rec, LnfFld.FldBrec1, brec)

	err = ptr.WriteRecord(&rec)

	if err != nil {
		fmt.Println(err)
	}
}
