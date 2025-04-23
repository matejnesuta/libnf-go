package record

import (
	"encoding/binary"
	"libnf/api/errors"
	"libnf/api/fields"
	"libnf/internal"
	"net"
	"strings"
	"time"
	"unsafe"
)

type Record struct {
	ptr       uintptr
	allocated bool
}

// GetPtr returns the pointer to the record.
func (r *Record) GetPtr() uintptr {
	return r.ptr
}

// Allocated returns whether the record is allocated.
func (r *Record) Allocated() bool {
	return r.allocated
}

func isAllBytesZero(data []byte) bool {
	n := len(data)

	// Round n down to the nearest multiple of 8
	// by clearing the last 3 bits.
	nlen8 := n & ^0b111
	i := 0

	for ; i < nlen8; i += 8 {
		b := *(*uint64)(unsafe.Pointer(&data[i]))
		if b != 0 {
			return false
		}
	}

	for ; i < n; i++ {
		if data[i] != 0 {
			return false
		}
	}

	return true
}

func convertIpToBytes(ip net.IP) []byte {
	newIP := ip.To4()
	if newIP != nil {
		return append(make([]byte, 12), []byte(newIP)...)
	}
	return []byte(ip.To16())
}

func convertToIP(data []byte) net.IP {
	if isAllBytesZero(data[:12]) {
		return net.IP(data[12:]).To4() // IPv4 Address
	}
	return net.IP(data) // IPv6 Address
}

func callFget(r *Record, field int, fieldPtr uintptr) error {
	status := internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
	if status == internal.ERR_NOTSET {
		return errors.ErrNotSet
	}
	return nil
}

func getSimpleDataType[T uint8 | uint16 | uint32 | uint64 | float64](r *Record, field int) (T, error) {
	var val T
	err := callFget(r, field, uintptr(unsafe.Pointer(&val)))
	if err != nil {
		return 0, err
	}
	return val, nil
}

func getIP(r *Record, field int) (any, error) {
	ipBuf := make([]byte, 16) // Allocate 16 bytes (IPv6 max size)
	fieldPtr := unsafe.Pointer(&ipBuf[0])
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return convertToIP(ipBuf), nil
}

func getTime(r *Record, field int) (any, error) {
	val := int64(0)
	fieldPtr := unsafe.Pointer(&val)
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return time.UnixMilli(val), nil

}

func getMacAddress(r *Record, field int) (any, error) {
	val := [6]byte{0}
	fieldPtr := unsafe.Pointer(&val)
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return net.HardwareAddr(val[:]), nil

}

func getBasicRecord1(r *Record, field int) (any, error) {
	brec := internal.NewLnf_brec1_t()
	defer internal.DeleteLnf_brec1_t(brec)
	err := callFget(r, field, uintptr(brec.Swigcptr()))

	if err != nil {
		return nil, err
	}

	var output fields.BasicRecord1
	output.First = time.Time(time.UnixMilli(int64(brec.GetFirst())))
	output.Last = time.Time(time.UnixMilli(int64(brec.GetLast())))
	output.Prot = brec.GetProt()
	output.SrcPort = brec.GetSrcport()
	output.DstPort = brec.GetDstport()
	output.Bytes = brec.GetBytes()
	output.Pkts = brec.GetPkts()
	output.Flows = brec.GetFlows()

	srcaddr := unsafe.Slice((*byte)(unsafe.Pointer(brec.GetSrcaddr().GetData())), 16)
	dstaddr := unsafe.Slice((*byte)(unsafe.Pointer(brec.GetDstaddr().GetData())), 16)

	output.SrcAddr = convertToIP(srcaddr)
	output.DstAddr = convertToIP(dstaddr)
	return output, nil
}

func getAcl(r *Record, field int) (any, error) {
	acl := internal.NewLnf_acl_t()
	defer internal.DeleteLnf_acl_t(acl)
	err := callFget(r, field, uintptr(acl.Swigcptr()))
	if err != nil {
		return nil, err
	}
	return fields.Acl{
		AclId:  acl.GetAcl_id(),
		AceId:  acl.GetAce_id(),
		XaceId: acl.GetXace_id(),
	}, nil
}

func getMpls(r *Record, field int) (any, error) {
	var mpls fields.Mpls
	mplsPtr := unsafe.Pointer(&mpls)
	err := callFget(r, field, uintptr(mplsPtr))
	if err != nil {
		return nil, err
	}
	return mpls, nil
}

// TODO test buffer ovewflow/underflow
func getString(r *Record, field int) (any, error) {
	buf := make([]byte, internal.STRING+1) // Adjust the size as needed
	fieldPtr := unsafe.Pointer(&buf[0])
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}

	return strings.TrimRight(string(buf), "\x00"), nil
}

func (r *Record) GetField(field int) (any, error) {
	if !r.allocated {
		return nil, errors.ErrRecordNotAllocated
	}

	expectedType, ok := fields.FieldTypes[field]
	var ret any
	var err error
	if !ok {
		return nil, errors.ErrUnknownFld
	}
	switch expectedType.(type) {
	case uint64:
		ret, err = getSimpleDataType[uint64](r, field)

	case uint32:
		ret, err = getSimpleDataType[uint32](r, field)

	case uint16:
		ret, err = getSimpleDataType[uint16](r, field)

	case uint8:
		ret, err = getSimpleDataType[uint8](r, field)

	case float64:
		ret, err = getSimpleDataType[float64](r, field)

	case net.IP:
		ret, err = getIP(r, field)

	case time.Time:
		ret, err = getTime(r, field)

	case net.HardwareAddr: // MacAddr
		ret, err = getMacAddress(r, field)

	case fields.BasicRecord1:
		ret, err = getBasicRecord1(r, field)

	case fields.Acl:
		ret, err = getAcl(r, field)

	case fields.Mpls:
		ret, err = getMpls(r, field)

	case string:
		ret, err = getString(r, field)

	default:
		ret = nil
		err = errors.ErrUnknownFld
	}
	return ret, err
}

func ipToUint32Array(ip net.IP) [4]uint32 {
	if ip == nil {
		return [4]uint32{}
	}

	// Ensure it's an IPv6 address
	if ip.To4() != nil {
		ip = net.IP(append(make([]byte, 12), ip...))
	}

	// Split the 16 bytes into 4 uint32 values
	var result [4]uint32
	for i := 0; i < 4; i++ {
		result[i] = binary.LittleEndian.Uint32(ip[i*4 : (i+1)*4])
	}

	return result
}

func SetField[T fields.FldDataType](r *Record, field int, value T) error {
	if !r.allocated {
		return errors.ErrRecordNotAllocated
	}

	expectedType, ok := fields.FieldTypes[field]
	if !ok {
		return errors.ErrUnknownFld
	}

	switch v := any(value).(type) {
	case uint64:
		_, ok := expectedType.(uint64)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		var val uint64 = v
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&val)))

	case uint32:
		_, ok := expectedType.(uint32)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		val := v
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&val)))

	case uint16:
		_, ok := expectedType.(uint16)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		val := v
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&val)))

	case uint8:
		_, ok := expectedType.(uint8)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		val := v
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&val)))

	case net.IP:
		_, ok := expectedType.(net.IP)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		addr := convertIpToBytes(v)
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&addr[0])))

	case time.Time:
		_, ok := expectedType.(time.Time)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		t := v.UnixMilli()
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&t)))

	case net.HardwareAddr:
		_, ok := expectedType.(net.HardwareAddr)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&v[0])))

	case fields.BasicRecord1:
		_, ok := expectedType.(fields.BasicRecord1)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		brec1 := internal.NewLnf_brec1_t()
		brec1.SetFirst(uint64(v.First.UnixMilli()))
		brec1.SetLast(uint64(v.Last.UnixMilli()))
		brec1.SetSrcport(v.SrcPort)
		brec1.SetDstport(v.DstPort)
		brec1.SetBytes(v.Bytes)
		brec1.SetPkts(v.Pkts)
		brec1.SetFlows(v.Flows)
		brec1.SetProt(v.Prot)

		srcaddr_t := internal.NewLnf_ip_t()
		dstaddr_t := internal.NewLnf_ip_t()

		srcip := ipToUint32Array(v.SrcAddr)
		dstip := ipToUint32Array(v.DstAddr)

		srcaddr_t.SetData(&srcip[0])
		dstaddr_t.SetData(&dstip[0])
		brec1.SetSrcaddr(srcaddr_t)
		brec1.SetDstaddr(dstaddr_t)

		internal.Rec_fset(r.ptr, field, uintptr(brec1.Swigcptr()))

		internal.DeleteLnf_brec1_t(brec1)
		internal.DeleteLnf_ip_t(srcaddr_t)
		internal.DeleteLnf_ip_t(dstaddr_t)

	case fields.Acl:
		_, ok := expectedType.(fields.Acl)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		aclPtr := internal.NewLnf_acl_t()
		aclPtr.SetAcl_id(v.AclId)
		aclPtr.SetAce_id(v.AceId)
		aclPtr.SetXace_id(v.XaceId)
		internal.Rec_fset(r.ptr, field, uintptr(aclPtr.Swigcptr()))
		internal.DeleteLnf_acl_t(aclPtr)

	case fields.Mpls:
		_, ok := expectedType.(fields.Mpls)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		mplsPtr := unsafe.Pointer(&v)
		internal.Rec_fset(r.ptr, field, uintptr(mplsPtr))

	case string:
		_, ok := expectedType.(string)
		if !ok {
			return errors.ErrMismatchingDataTypes
		}
		buf := append([]byte(v), 0)
		internal.Rec_fset(r.ptr, field, uintptr(unsafe.Pointer(&buf[0])))

	default:
		return errors.ErrUnknownFld
	}

	return nil
}

func NewRecord() (Record, error) {
	var r Record
	status := internal.Rec_init(&r.ptr)
	if status == internal.ERR_NOMEM {
		return r, errors.ErrNoMem
	} else if status == internal.ERR_OTHER {
		return r, errors.ErrOther
	}
	r.allocated = true
	return r, nil
}

func (r *Record) Free() error {
	if r.allocated {
		internal.Rec_free(r.ptr)
		r.allocated = false
		return nil
	}
	return errors.ErrRecordNotAllocated
}
func (r *Record) Clear() error {
	if r.allocated {
		internal.Rec_clear(r.ptr)
		return nil
	}
	return errors.ErrRecordNotAllocated
}

func (r *Record) CopyFrom(other Record) error {
	if !r.allocated || !other.allocated {
		return errors.ErrRecordNotAllocated
	}

	status := internal.Rec_copy(r.ptr, other.ptr)
	if status == internal.ERR_OTHER {
		return errors.ErrOther
	}
	return nil
}
