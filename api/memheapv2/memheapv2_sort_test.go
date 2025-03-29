package memheapv2_test

import (
	"fmt"
	"libnf/api/errors"
	"libnf/api/fields"
	memheap "libnf/api/memheapv2"
	"libnf/api/record"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sortByUint8(t *testing.T, ports []uint16, protocols []uint8, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Prot, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [3]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(111),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("1.1.1.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.3").To4(),
		DstAddr: net.ParseIP("1.1.1.4").To4(),
		Prot:    uint8(5),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(90),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.5").To4(),
		DstAddr: net.ParseIP("1.1.1.6").To4(),
		Prot:    uint8(4),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, protocols[i], brec.Prot)
		i++
	}
}

func sortByUint16(t *testing.T, bytes []uint64, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.Doctets, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [3]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(111),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("1.1.1.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.3").To4(),
		DstAddr: net.ParseIP("1.1.1.4").To4(),
		Prot:    uint8(5),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(90),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.5").To4(),
		DstAddr: net.ParseIP("1.1.1.6").To4(),
		Prot:    uint8(4),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
}

func sortByUint32(t *testing.T, as []uint32, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.EgressAclId, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()
	var inputAs = [3]uint32{2, 3, 1}
	var inputSrcPort = [3]uint16{80, 443, 53}
	for i := 0; i < 3; i++ {

		err = record.SetField(&rec, fields.EgressAclId, inputAs[i])
		assert.Nil(t, err)
		err = record.SetField(&rec, fields.SrcPort, inputSrcPort[i])
		assert.Nil(t, err)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.EgressAclId)
		assert.Equal(t, as[i], val)
		val, _ = rec.GetField(fields.SrcPort)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func sortByUint64(t *testing.T, bytes []uint64, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.Doctets, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [3]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(111),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("1.1.1.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(70),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.3").To4(),
		DstAddr: net.ParseIP("1.1.1.4").To4(),
		Prot:    uint8(5),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(90),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("1.1.1.5").To4(),
		DstAddr: net.ParseIP("1.1.1.6").To4(),
		Prot:    uint8(4),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		fmt.Println(brec.Bytes, brec.SrcPort)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
}

func sortByTime(t *testing.T, times []time.Time, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.First, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	record.SetField(&rec, fields.First, time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local))
	record.SetField(&rec, fields.SrcPort, uint16(111))
	err = heap.WriteRecord(&rec)
	assert.Nil(t, err)

	record.SetField(&rec, fields.First, time.Date(2018, time.May, 28, 15, 55, 0, 0, time.Local))
	record.SetField(&rec, fields.SrcPort, uint16(80))
	err = heap.WriteRecord(&rec)
	assert.Nil(t, err)

	record.SetField(&rec, fields.First, time.Date(2015, time.May, 28, 15, 55, 0, 0, time.Local))
	record.SetField(&rec, fields.SrcPort, uint16(90))
	err = heap.WriteRecord(&rec)
	assert.Nil(t, err)

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.First)
		assert.Equal(t, times[i], val)
		val, _ = rec.GetField(fields.SrcPort)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func sortByIP(t *testing.T, srcAddrs []net.IP, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.SrcAddr, memheap.AggrKey, order, 32, 128)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	inputIps := [4]net.IP{net.ParseIP("192.168.1.1").To4(),
		net.ParseIP("192.168.1.2").To4(),
		net.ParseIP("2001:399a:eddf:f709:ff00:ff00:2ce5:7918"),
		net.ParseIP("2001:399a:eddf:f709:ff00:ff00:2ce5:7917")}
	inputPorts := [4]uint16{80, 443, 53, 53}

	for i := 0; i < 4; i++ {
		fmt.Println(inputIps[i], inputPorts[i])
		record.SetField(&rec, fields.SrcAddr, inputIps[i])
		record.SetField(&rec, fields.SrcPort, inputPorts[i])
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 4 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.SrcAddr)
		assert.Equal(t, srcAddrs[i], val)
		fmt.Print(val, " ")
		val, _ = rec.GetField(fields.SrcPort)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func sortByMac(t *testing.T, srcMacs []net.HardwareAddr, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.InSrcMac, memheap.AggrKey, order, 48, 128)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	inputMacs := [3]net.HardwareAddr{net.HardwareAddr{0x00, 0x02, 0x00, 0x00, 0x00, 0x00},
		net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x02}}
	inputPorts := [3]uint16{80, 443, 53}

	for i := 0; i < 3; i++ {
		record.SetField(&rec, fields.InSrcMac, inputMacs[i])
		record.SetField(&rec, fields.SrcPort, inputPorts[i])
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.InSrcMac)
		assert.Equal(t, srcMacs[i], val)
		fmt.Print(val, " ")
		val, _ = rec.GetField(fields.SrcPort)
		fmt.Println(val)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func sortByDuration(t *testing.T, durations []uint64, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.CalcDuration, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	inputDurations := [3]uint64{70, 20, 80}
	inputPorts := [3]uint16{80, 443, 53}

	for i := 0; i < 3; i++ {
		record.SetField(&rec, fields.First, time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local))
		record.SetField(&rec, fields.Last, time.Date(2017, time.May, 28, 15, 55, 0, int(inputDurations[i]*1000000), time.Local))
		record.SetField(&rec, fields.SrcPort, inputPorts[i])
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.CalcDuration)
		assert.Equal(t, durations[i], val)
		val, _ = rec.GetField(fields.SrcPort)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func sortByBps(t *testing.T, bpss []float64, ports []uint16, order int, calcField int, item int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(calcField, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	inputBytes := [3]uint64{70, 20, 80}
	inputPorts := [3]uint16{80, 443, 53}

	for i := 0; i < 3; i++ {
		record.SetField(&rec, fields.First, time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local))
		record.SetField(&rec, fields.Last, time.Date(2017, time.May, 28, 15, 55, 2, 0, time.Local))
		record.SetField(&rec, item, uint64(inputBytes[i]))
		record.SetField(&rec, fields.SrcPort, inputPorts[i])
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(calcField)
		fmt.Print(val, " ")
		assert.Equal(t, bpss[i], val)
		val, _ = rec.GetField(fields.SrcPort)
		fmt.Println(val)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func sortByBpp(t *testing.T, bpps []float64, ports []uint16, order int) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.CalcBpp, memheap.AggrAuto, order, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.SrcPort, memheap.AggrAuto, memheap.SortNone, 0, 0)
	assert.Nil(t, err)

	rec, _ := record.NewRecord()
	defer rec.Free()

	inputBytes := [3]uint64{70, 20, 80}
	inputPorts := [3]uint16{80, 443, 53}

	for i := 0; i < 3; i++ {
		record.SetField(&rec, fields.Dpkts, uint64(2))
		record.SetField(&rec, fields.SrcPort, inputPorts[i])
		record.SetField(&rec, fields.Doctets, uint64(inputBytes[i]))
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.CalcBpp)
		fmt.Print(val, " ")
		assert.Equal(t, bpps[i], val)
		val, _ = rec.GetField(fields.SrcPort)
		fmt.Println(val)
		assert.Equal(t, ports[i], val)
		i++
	}
}

func TestSortByUint8Asc(t *testing.T) {
	ports := [3]uint16{90, 80, 111}
	protocols := [3]uint8{4, 5, 6}
	sortByUint8(t, ports[:], protocols[:], memheap.SortAsc)
}

func TestSortByUint8Desc(t *testing.T) {
	ports := [3]uint16{111, 80, 90}
	protocols := [3]uint8{6, 5, 4}
	sortByUint8(t, ports[:], protocols[:], memheap.SortDesc)
}

func TestSortByUint16Asc(t *testing.T) {
	bytes := [3]uint64{80, 80, 20}
	ports := [3]uint16{80, 90, 111}
	sortByUint16(t, bytes[:], ports[:], memheap.SortAsc)
}

func TestSortByUint16Desc(t *testing.T) {
	bytes := [3]uint64{20, 80, 80}
	ports := [3]uint16{111, 90, 80}
	sortByUint16(t, bytes[:], ports[:], memheap.SortDesc)
}

func TestSortByUint32Asc(t *testing.T) {
	as := [3]uint32{1, 2, 3}
	ports := [3]uint16{53, 80, 443}
	sortByUint32(t, as[:], ports[:], memheap.SortAsc)
}

func TestSortByUint32Desc(t *testing.T) {
	as := [3]uint32{3, 2, 1}
	ports := [3]uint16{443, 80, 53}
	sortByUint32(t, as[:], ports[:], memheap.SortDesc)
}

func TestSortByUint64Asc(t *testing.T) {
	bytes := [3]uint64{20, 70, 80}
	ports := [3]uint16{111, 80, 90}
	sortByUint64(t, bytes[:], ports[:], memheap.SortAsc)
}

func TestSortByUint64Desc(t *testing.T) {
	bytes := [3]uint64{80, 70, 20}
	ports := [3]uint16{90, 80, 111}
	sortByUint64(t, bytes[:], ports[:], memheap.SortDesc)
}

func TestSortByTimeAsc(t *testing.T) {
	times := [3]time.Time{
		time.Date(2015, time.May, 28, 15, 55, 0, 0, time.Local),
		time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local),
		time.Date(2018, time.May, 28, 15, 55, 0, 0, time.Local),
	}
	ports := [3]uint16{90, 111, 80}
	sortByTime(t, times[:], ports[:], memheap.SortAsc)
}

func TestSortByTimeDesc(t *testing.T) {
	times := [3]time.Time{
		time.Date(2018, time.May, 28, 15, 55, 0, 0, time.Local),
		time.Date(2017, time.May, 28, 15, 55, 0, 0, time.Local),
		time.Date(2015, time.May, 28, 15, 55, 0, 0, time.Local),
	}
	ports := [3]uint16{80, 111, 90}
	sortByTime(t, times[:], ports[:], memheap.SortDesc)
}

// these 2 tests might seem wrong, but this is how libnf sorts IP addresses at the moment
func TestSortByIPAsc(t *testing.T) {
	ips := [4]net.IP{
		net.ParseIP("2001:399a:eddf:f709:ff00:ff00:2ce5:7918"),
		net.ParseIP("2001:399a:eddf:f709:ff00:ff00:2ce5:7917"),
		net.ParseIP("192.168.1.2").To4(),
		net.ParseIP("192.168.1.1").To4()}

	ports := [4]uint16{53, 53, 443, 80}
	sortByIP(t, ips[:], ports[:], memheap.SortAsc)
}

func TestSortByIPDesc(t *testing.T) {
	ips := [4]net.IP{
		net.ParseIP("192.168.1.1").To4(),
		net.ParseIP("192.168.1.2").To4(),
		net.ParseIP("2001:399a:eddf:f709:ff00:ff00:2ce5:7917"),
		net.ParseIP("2001:399a:eddf:f709:ff00:ff00:2ce5:7918")}

	ports := [4]uint16{80, 443, 53, 53}
	sortByIP(t, ips[:], ports[:], memheap.SortDesc)
}

func TestSortByMacAsc(t *testing.T) {
	macs := [3]net.HardwareAddr{
		net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		net.HardwareAddr{0x00, 0x02, 0x00, 0x00, 0x00, 0x00}}
	ports := [3]uint16{443, 53, 80}
	sortByMac(t, macs[:], ports[:], memheap.SortAsc)
}

func TestSortByMacDesc(t *testing.T) {
	macs := [3]net.HardwareAddr{
		net.HardwareAddr{0x00, 0x02, 0x00, 0x00, 0x00, 0x00},
		net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}
	ports := [3]uint16{80, 53, 443}
	sortByMac(t, macs[:], ports[:], memheap.SortDesc)
}

func TestSortByDurationAsc(t *testing.T) {
	durations := [3]uint64{20, 70, 80}
	ports := [3]uint16{443, 80, 53}
	sortByDuration(t, durations[:], ports[:], memheap.SortAsc)
}

func TestSortByDurationDesc(t *testing.T) {
	durations := [3]uint64{80, 70, 20}
	ports := [3]uint16{53, 80, 443}
	sortByDuration(t, durations[:], ports[:], memheap.SortDesc)
}

func TestSortByBpsAsc(t *testing.T) {
	bpss := [3]float64{80, 280, 320}
	ports := [3]uint16{443, 80, 53}
	sortByBps(t, bpss[:], ports[:], memheap.SortAsc, fields.CalcBps, fields.Doctets)
}

func TestSortByBpsDesc(t *testing.T) {
	bpss := [3]float64{320, 280, 80}
	ports := [3]uint16{53, 80, 443}
	sortByBps(t, bpss[:], ports[:], memheap.SortDesc, fields.CalcBps, fields.Doctets)
}

func TestSortByPpsAsc(t *testing.T) {
	bpss := [3]float64{10, 35, 40}
	ports := [3]uint16{443, 80, 53}
	sortByBps(t, bpss[:], ports[:], memheap.SortAsc, fields.CalcPps, fields.Dpkts)
}

func TestSortByPpsDesc(t *testing.T) {
	bpss := [3]float64{40, 35, 10}
	ports := [3]uint16{53, 80, 443}
	sortByBps(t, bpss[:], ports[:], memheap.SortDesc, fields.CalcPps, fields.Dpkts)
}

func TestSortByBppAsc(t *testing.T) {
	bpss := [3]float64{10, 35, 40}
	ports := [3]uint16{443, 80, 53}
	sortByBpp(t, bpss[:], ports[:], memheap.SortAsc)
}

func TestSortByBppDesc(t *testing.T) {
	bpss := [3]float64{40, 35, 10}
	ports := [3]uint16{53, 80, 443}
	sortByBpp(t, bpss[:], ports[:], memheap.SortDesc)
}

func TestAggrPerPairField(t *testing.T) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SetNfdumpComp(false)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [2]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(53),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("2.2.2.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("3.3.3.3").To4(),
		DstAddr: net.ParseIP("4.4.4.4").To4(),
		Prot:    uint8(6),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	ports := [3]uint16{53, 80, 1222}
	pkts := [3]uint64{2, 3, 3}
	bytes := [3]uint64{40, 80, 80}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for i < 3 {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, pkts[i], brec.Pkts)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
}

func TestAggrPerPairFieldWithNfdumpComp(t *testing.T) {
	var heap memheap.MemHeapV2 = *memheap.NewMemHeapV2()
	err := heap.SortAggrOptions(fields.PairPort, memheap.AggrKey, memheap.SortAsc, 0, 0)
	assert.Nil(t, err)
	err = heap.SortAggrOptions(fields.Dpkts, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SortAggrOptions(fields.Doctets, memheap.AggrSum, memheap.SortNone, 0, 0)
	assert.Nil(t, err)
	heap.SetNfdumpComp(true)

	rec, _ := record.NewRecord()
	defer rec.Free()

	brecs := [2]fields.BasicRecord1{{
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(20),
		Pkts:    uint64(1),
		Flows:   uint64(1),
		SrcPort: uint16(53),
		DstPort: uint16(53),
		SrcAddr: net.ParseIP("1.1.1.1").To4(),
		DstAddr: net.ParseIP("2.2.2.2").To4(),
		Prot:    uint8(6),
	}, {
		First:   time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Last:    time.Date(2017, time.May, 28, 15, 55, 0, 0, time.UTC),
		Bytes:   uint64(80),
		Pkts:    uint64(3),
		Flows:   uint64(1),
		SrcPort: uint16(80),
		DstPort: uint16(1222),
		SrcAddr: net.ParseIP("3.3.3.3").To4(),
		DstAddr: net.ParseIP("4.4.4.4").To4(),
		Prot:    uint8(6),
	}}

	for _, brec := range brecs {
		record.SetField(&rec, fields.Brec1, brec)
		err := heap.WriteRecord(&rec)
		assert.Equal(t, nil, err)
	}

	ports := [3]uint16{53, 80, 1222}
	pkts := [3]uint64{1, 3, 3}
	bytes := [3]uint64{20, 80, 80}

	i := 0
	cursor, err := heap.FirstRecordPosition()
	assert.Nil(t, err)
	for {
		err := heap.GetRecord(&cursor, &rec)
		assert.Nil(t, err)
		cursor, err = heap.NextRecordPosition(cursor)
		if err == errors.ErrMemHeapEnd {
			break
		}
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Brec1)
		brec := val.(fields.BasicRecord1)
		assert.Equal(t, ports[i], brec.SrcPort)
		assert.Equal(t, pkts[i], brec.Pkts)
		assert.Equal(t, bytes[i], brec.Bytes)
		i++
	}
}
