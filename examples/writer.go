package examples

import (
	"fmt"
	libnf "libnf/api"
	"net"
	"time"
)

// Writer is a function that demonstrates how to use the libnf package to write data to a file.
func Writer() {
	var ptr libnf.File
	err := ptr.OpenWrite("api/tests/testfiles/test-file.tmp", "tmp/output.tmp", false, 0, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := libnf.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	// Set the fields of the record

	t := time.Now()
	brec := libnf.BasicRecord1{
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

	acl := libnf.Acl{
		AclId:  1,
		AceId:  2,
		XaceId: 3,
	}

	mpls := libnf.Mpls{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	username := "leopold kasing"

	mac := net.HardwareAddr{0x00, 0x08, 0x00, 0x0c, 0x02, 0x01}

	exporterIp := net.IPv4(10, 10, 0, 1)
	var exporterId uint32 = 10

	var bps float64 = 1000

	libnf.SetField(&rec, libnf.FldCalcBps, bps)
	libnf.SetField(&rec, libnf.FldSrcAS, exporterId)
	libnf.SetField(&rec, libnf.FldBrec1, brec)
	libnf.SetField(&rec, libnf.FldIngressAcl, acl)
	libnf.SetField(&rec, libnf.FldMplsLabel, mpls)
	libnf.SetField(&rec, libnf.FldUsername, username)
	libnf.SetField(&rec, libnf.FldInSrcMac, mac)

	libnf.SetField(&rec, libnf.FldIpNextHop, exporterIp)

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
	libnf.SetField(&rec, libnf.FldExporterIp, net.IP(exporterIpBytes[:]))
	libnf.SetField(&rec, libnf.FldBrec1, brec)

	err = ptr.WriteRecord(&rec)

	if err != nil {
		fmt.Println(err)
	}
}
