package record_test

import (
	"fmt"
	LnfErr "libnf/api/errors"
	LnfFld "libnf/api/fields"
	LnfFile "libnf/api/file"
	LnfRec "libnf/api/record"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecordInit(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	assert.Equal(t, rec.Allocated(), true)
	err = rec.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, rec.Allocated(), false)
}

func TestFreeUnallocatedRecord(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	err = rec.Free()
	assert.Equal(t, nil, err)
	err = rec.Free()
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}

func TestGetBrec1FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldBrec1)
	brec, ok := val.(LnfFld.BasicRecord1)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(6), brec.Prot)
	assert.Equal(t, uint16(1123), brec.SrcPort)
	assert.Equal(t, uint16(80), brec.DstPort)
	assert.Equal(t, uint64(12345), brec.Bytes)
	assert.Equal(t, time.UnixMilli(11220), brec.First)
	assert.Equal(t, time.UnixMilli(11229), brec.Last)
	assert.Equal(t, uint64(20), brec.Pkts)
	assert.Equal(t, uint64(1), brec.Flows)
	assert.Equal(t, net.IPv4(192, 168, 0, 1).To4(), brec.SrcAddr)
	assert.Equal(t, net.IPv4(192, 168, 0, 2).To4(), brec.DstAddr)
}

func TestGetBrec1WithIPv6FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv6-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldBrec1)
	brec, ok := val.(LnfFld.BasicRecord1)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(6), brec.Prot)
	assert.Equal(t, uint16(1123), brec.SrcPort)
	assert.Equal(t, uint16(80), brec.DstPort)
	assert.Equal(t, uint64(12345), brec.Bytes)
	assert.Equal(t, time.UnixMilli(11220), brec.First)
	assert.Equal(t, time.UnixMilli(11229), brec.Last)
	assert.Equal(t, uint64(20), brec.Pkts)
	assert.Equal(t, uint64(1), brec.Flows)
	assert.Equal(t, net.ParseIP("2001:67c:1220::aa:bb"), brec.SrcAddr)
	assert.Equal(t, net.ParseIP("2001:67c:1220::11:22"), brec.DstAddr)
}

func TestGetAclFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldIngressAcl)
	acl, ok := val.(LnfFld.Acl)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint32(25), acl.AclId)
	assert.Equal(t, uint32(25), acl.AceId)
	assert.Equal(t, uint32(25), acl.XaceId)
}

func TestGetMplsFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldMplsLabel)
	mpls, ok := val.(LnfFld.Mpls)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint32(10), mpls[0])
	assert.Equal(t, uint32(9), mpls[1])
	assert.Equal(t, uint32(8), mpls[2])
	assert.Equal(t, uint32(7), mpls[3])
	assert.Equal(t, uint32(6), mpls[4])
	assert.Equal(t, uint32(5), mpls[5])
	assert.Equal(t, uint32(4), mpls[6])
	assert.Equal(t, uint32(3), mpls[7])
	assert.Equal(t, uint32(2), mpls[8])
	assert.Equal(t, uint32(1), mpls[9])
}

func TestGetMacFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldInSrcMac)
	mac, ok := val.(net.HardwareAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, mac)
}

func TestGetUsernameFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldUsername)
	username, ok := val.(string)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, "testuser", username)
}

func TestGetIpFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldSrcaddr)
	ip, ok := val.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, net.IPv4(192, 168, 0, 1).To4(), ip)
}

func TestGetUint32FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldIngressAclId)
	id, ok := val.(uint32)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint32(25), id)
}

func TestGetUint16FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldSrcport)
	port, ok := val.(uint16)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint16(1123), port)
}

func TestGetUint8FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldProt)
	prot, ok := val.(uint8)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(6), prot)
}

func TestGetUint64FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldDoctets)
	doctets, ok := val.(uint64)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint64(12345), doctets)
}

func TestGetFloat64FromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv6-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldCalcBps)
	bps, ok := val.(float64)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, float64(1.0973333333333334e+07), bps)
}

func TestGetTimeFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv6-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldFirst)
	first, ok := val.(time.Time)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, time.UnixMilli(11220), first)
}

func TestGetUnsetFieldFromRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldIpNextHop)
	assert.Equal(t, LnfErr.ErrNotSet, err)
	assert.Equal(t, nil, val)
}

func TestGetFieldFromUnallocatedRecord(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	err = rec.Free()
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldIpNextHop)
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
	assert.Equal(t, nil, val)
}

func TestGetNonexistentField(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec, err := LnfRec.NewRecord()
	defer rec.Free()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec)

	val, err := rec.GetField(999)
	assert.Equal(t, LnfErr.ErrUnknownFld, err)
	assert.Equal(t, nil, val)
}

func TestCopyRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()

	rec1, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	rec2, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	ptr.GetNextRecord(&rec1)
	defer rec1.Free()
	defer rec2.Free()

	err = rec1.CopyFrom(rec2)
	assert.Equal(t, nil, err)

	val1, err := rec1.GetField(LnfFld.FldSrcaddr)
	ip1, ok := val1.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	val2, err := rec2.GetField(LnfFld.FldSrcaddr)
	ip2, ok := val2.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, ip1, ip2)
}

func TestCopyFromUnallocatedRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()
	rec1, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	rec2, err := LnfRec.NewRecord()
	ptr.GetNextRecord(&rec2)
	defer rec1.Free()

	assert.Equal(t, nil, err)
	err = rec2.Free()
	assert.Equal(t, nil, err)

	err = rec1.CopyFrom(rec2)
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}

func TestCopyToUnallocatedRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()
	rec1, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	rec2, err := LnfRec.NewRecord()
	ptr.GetNextRecord(&rec2)
	defer rec2.Free()

	assert.Equal(t, nil, err)
	err = rec1.Free()
	assert.Equal(t, nil, err)

	err = rec1.CopyFrom(rec2)
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}

func TestClearRecord(t *testing.T) {
	var ptr LnfFile.File
	err := ptr.OpenRead("../testfiles/ipv4-file.tmp", false, false)
	assert.Equal(t, nil, err)
	defer ptr.Close()
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	ptr.GetNextRecord(&rec)
	val, err := rec.GetField(LnfFld.FldSrcaddr)
	ip, ok := val.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)
	assert.Equal(t, net.IPv4(192, 168, 0, 1).To4(), ip)

	err = rec.Clear()
	assert.Equal(t, nil, err)

	val, err = rec.GetField(LnfFld.FldSrcaddr)
	assert.Equal(t, nil, err)
	ip, ok = val.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), ip)

	assert.Equal(t, nil, err)
}

func TestClearUnallocatedRecord(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	err = rec.Free()
	assert.Equal(t, nil, err)

	err = rec.Clear()
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}

func TestSetFieldUint32(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldIngressAclId, uint32(25))
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldIngressAclId)
	id, ok := val.(uint32)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint32(25), id)
}

func TestSetFieldUint16(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldSrcport, uint16(1123))
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldSrcport)
	port, ok := val.(uint16)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint16(1123), port)
}

func TestSetFieldUint8(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldProt, uint8(6))
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldProt)
	prot, ok := val.(uint8)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(6), prot)
}

func TestSetFieldUint64(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldDoctets, uint64(12345))
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldDoctets)
	doctets, ok := val.(uint64)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint64(12345), doctets)
}

func TestSetFieldTime(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldFirst, time.UnixMilli(11220))
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldFirst)
	first, ok := val.(time.Time)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, time.UnixMilli(11220), first)
}

func TestSetFieldMpls(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldMplsLabel, LnfFld.Mpls{10, 9, 8, 7, 6, 5, 4, 3, 2, 1})
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldMplsLabel)
	mpls, ok := val.(LnfFld.Mpls)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, LnfFld.Mpls{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, mpls)
}

func TestSetFieldMac(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldInSrcMac, net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldInSrcMac)
	mac, ok := val.(net.HardwareAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, mac)
}

func TestSetFieldUsername(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()

	err = LnfRec.SetField(&rec, LnfFld.FldUsername, "testuser")
	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldUsername)
	username, ok := val.(string)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, "testuser", username)
}

func TestSetFieldIp(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	err = LnfRec.SetField(&rec, LnfFld.FldSrcaddr, net.IPv4(192, 168, 0, 1).To4())

	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldSrcaddr)
	ip, ok := val.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, net.IPv4(192, 168, 0, 1).To4(), ip)
}

func TestSetFieldIpv6(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	err = LnfRec.SetField(&rec, LnfFld.FldSrcaddr, net.ParseIP("2001:67c:1220::aa:bb"))

	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldSrcaddr)
	ip, ok := val.(net.IP)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, net.ParseIP("2001:67c:1220::aa:bb"), ip)
}

func TestSetFieldAcl(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	err = LnfRec.SetField(&rec, LnfFld.FldIngressAcl, LnfFld.Acl{25, 25, 25})

	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldIngressAcl)
	acl, ok := val.(LnfFld.Acl)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, LnfFld.Acl{25, 25, 25}, acl)
}

func TestSetFieldBrec1(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	input := LnfFld.BasicRecord1{
		Prot:    6,
		SrcPort: 1123,
		DstPort: 80,
		Bytes:   12345,
		First:   time.UnixMilli(11220),
		Last:    time.UnixMilli(11229),
		Pkts:    20,
		Flows:   1,
		SrcAddr: net.IPv4(192, 168, 0, 1).To4(),
		DstAddr: net.IPv4(192, 168, 0, 2).To4(),
	}

	fmt.Println([]byte(input.SrcAddr))
	err = LnfRec.SetField(&rec, LnfFld.FldBrec1, input)

	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldBrec1)
	brec, ok := val.(LnfFld.BasicRecord1)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(6), brec.Prot)
	assert.Equal(t, uint16(1123), brec.SrcPort)
	assert.Equal(t, uint16(80), brec.DstPort)
	assert.Equal(t, uint64(12345), brec.Bytes)
	assert.Equal(t, time.UnixMilli(11220), brec.First)
	assert.Equal(t, time.UnixMilli(11229), brec.Last)
	assert.Equal(t, uint64(20), brec.Pkts)
	assert.Equal(t, uint64(1), brec.Flows)
	assert.Equal(t, net.IPv4(192, 168, 0, 1).To4(), brec.SrcAddr)
	assert.Equal(t, net.IPv4(192, 168, 0, 2).To4(), brec.DstAddr)
}

func TestSetFieldBrec1Ipv6(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	err = LnfRec.SetField(&rec, LnfFld.FldBrec1, LnfFld.BasicRecord1{
		Prot:    6,
		SrcPort: 1123,
		DstPort: 80,
		Bytes:   12345,
		First:   time.UnixMilli(11220),
		Last:    time.UnixMilli(11229),
		Pkts:    20,
		Flows:   1,
		SrcAddr: net.ParseIP("2001:67c:1220::aa:bb"),
		DstAddr: net.ParseIP("2001:67c:1220::11:22"),
	})

	assert.Equal(t, nil, err)

	val, err := rec.GetField(LnfFld.FldBrec1)
	brec, ok := val.(LnfFld.BasicRecord1)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(6), brec.Prot)
	assert.Equal(t, uint16(1123), brec.SrcPort)
	assert.Equal(t, uint16(80), brec.DstPort)
	assert.Equal(t, uint64(12345), brec.Bytes)
	assert.Equal(t, time.UnixMilli(11220), brec.First)
	assert.Equal(t, time.UnixMilli(11229), brec.Last)
	assert.Equal(t, uint64(20), brec.Pkts)
	assert.Equal(t, uint64(1), brec.Flows)
	assert.Equal(t, net.ParseIP("2001:67c:1220::aa:bb"), brec.SrcAddr)
	assert.Equal(t, net.ParseIP("2001:67c:1220::11:22"), brec.DstAddr)
}

func TestSetFieldUnknownField(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	err = LnfRec.SetField(&rec, 999, "test")

	assert.Equal(t, LnfErr.ErrUnknownFld, err)
}

func TestSetFieldMismatchingDataTypes(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	defer rec.Free()
	err = LnfRec.SetField(&rec, LnfFld.FldSrcaddr, uint32(25))

	assert.Equal(t, LnfErr.ErrMismatchingDataTypes, err)
}

func TestSetFieldUnallocatedRecord(t *testing.T) {
	rec, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	err = rec.Free()
	assert.Equal(t, nil, err)

	err = LnfRec.SetField(&rec, LnfFld.FldSrcaddr, net.IPv4(192, 168, 0, 1).To4())
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}
