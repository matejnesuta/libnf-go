package libnf

import (
	"libnf/internal"
	"net"
	"reflect"
	"time"
	"unsafe"
)

type Record struct {
	ptr       uintptr
	allocated bool
}

type MacAddr net.HardwareAddr

type BasicRecord1 struct {
	First   time.Time // LNF_FLD_FIRST
	Last    time.Time // LNF_FLD_LAST
	SrcAddr net.IP    // LNF_FLD_SRCADDR
	DstAddr net.IP    // LNF_FLD_DSTADDR
	Prot    uint8     // LNF_FLD_PROT
	SrcPort uint16    // LNF_FLD_SRCPORT
	DstPort uint16    // LNF_FLD_DSTPORT
	Bytes   uint64    // LNF_FLD_DOCTETS
	Pkts    uint64    // LNF_FLD_DPKTS
	Flows   uint64    // LNF_FLD_AGGR_FLOWS
}

type Acl struct {
	AclId  uint32
	AceId  uint32
	XaceId uint32
}

const (
	FldFirst             int = internal.FLD_FIRST
	FldLast              int = internal.FLD_LAST
	FldReceived          int = internal.FLD_RECEIVED
	FldDoctets           int = internal.FLD_DOCTETS
	FldDpkts             int = internal.FLD_DPKTS
	FldDpktsAlias        int = internal.FLD_DPKTS_ALIAS
	FldOutBytes          int = internal.FLD_OUT_BYTES
	FldOutPkts           int = internal.FLD_OUT_PKTS
	FldOutPktsAlias      int = internal.FLD_OUT_PKTS_ALIAS
	FldAggrFlows         int = internal.FLD_AGGR_FLOWS
	FldSrcport           int = internal.FLD_SRCPORT
	FldDstport           int = internal.FLD_DSTPORT
	FldTcpFlags          int = internal.FLD_TCP_FLAGS
	FldTcpFlagsAlias     int = internal.FLD_TCP_FLAGS_ALIAS
	FldSrcaddr           int = internal.FLD_SRCADDR
	FldDstaddr           int = internal.FLD_DSTADDR
	FldSrcaddrAlias      int = internal.FLD_SRCADDR_ALIAS
	FldDstaddrAlias      int = internal.FLD_DSTADDR_ALIAS
	FldIpNextHop         int = internal.FLD_IP_NEXTHOP
	FldIpNextHopAlias    int = internal.FLD_IP_NEXTHOP_ALIAS
	FldSrcMask           int = internal.FLD_SRC_MASK
	FldDstMask           int = internal.FLD_DST_MASK
	FldTos               int = internal.FLD_TOS
	FldDstTos            int = internal.FLD_DST_TOS
	FldSrcAS             int = internal.FLD_SRCAS
	FldDstAS             int = internal.FLD_DSTAS
	FldBgpNextAdjacentAS int = internal.FLD_BGPNEXTADJACENTAS
	FldBgpPrevAdjacentAS int = internal.FLD_BGPPREVADJACENTAS
	FldBgpNextHop        int = internal.FLD_BGP_NEXTHOP
	FldProt              int = internal.FLD_PROT
	FldSrcVlan           int = internal.FLD_SRC_VLAN
	FldDstVlan           int = internal.FLD_DST_VLAN
	FldInSrcMac          int = internal.FLD_IN_SRC_MAC
	FldOutSrcMac         int = internal.FLD_OUT_SRC_MAC
	FldInDstMac          int = internal.FLD_IN_DST_MAC
	FldOutDstMac         int = internal.FLD_OUT_DST_MAC
	FldMplsLabel         int = internal.FLD_MPLS_LABEL
	FldInput             int = internal.FLD_INPUT
	FldOutput            int = internal.FLD_OUTPUT
	FldDir               int = internal.FLD_DIR
	FldFwdStatus         int = internal.FLD_FWD_STATUS
	FldIpRouter          int = internal.FLD_IP_ROUTER
	FldIpRouterAlias     int = internal.FLD_IP_ROUTER_ALIAS
	FldEngineType        int = internal.FLD_ENGINE_TYPE
	FldEngineId          int = internal.FLD_ENGINE_ID
	FldEngineTypeAlias   int = internal.FLD_ENGINE_TYPE_ALIAS
	FldEngineIdAlias     int = internal.FLD_ENGINE_ID_ALIAS
	FldEventTime         int = internal.FLD_EVENT_TIME
	FldConnId            int = internal.FLD_CONN_ID
	FldIcmpCode          int = internal.FLD_ICMP_CODE
	FldIcmpType          int = internal.FLD_ICMP_TYPE
	FldIcmpCodeAlias     int = internal.FLD_ICMP_CODE_ALIAS
	FldIcmpTypeAlias     int = internal.FLD_ICMP_TYPE_ALIAS
	FldFwXEvent          int = internal.FLD_FW_XEVENT
	FldFwEvent           int = internal.FLD_FW_EVENT
	FldXlateSrcIp        int = internal.FLD_XLATE_SRC_IP
	FldXlateDstIp        int = internal.FLD_XLATE_DST_IP
	FldXlateSrcPort      int = internal.FLD_XLATE_SRC_PORT
	FldXlateDstPort      int = internal.FLD_XLATE_DST_PORT
	FldIngressAclId      int = internal.FLD_INGRESS_ACL_ID
	FldIngressAceId      int = internal.FLD_INGRESS_ACE_ID
	FldIngressXaceId     int = internal.FLD_INGRESS_XACE_ID
	FldIngressAcl        int = internal.FLD_INGRESS_ACL
	FldEgressAclId       int = internal.FLD_EGRESS_ACL_ID
	FldEgressAceId       int = internal.FLD_EGRESS_ACE_ID
	FldEgressXaceId      int = internal.FLD_EGRESS_XACE_ID
	FldEgressAcl         int = internal.FLD_EGRESS_ACL
	FldUsername          int = internal.FLD_USERNAME
	FldIngressVrfid      int = internal.FLD_INGRESS_VRFID
	FldEventFlag         int = internal.FLD_EVENT_FLAG
	FldEgressVrfid       int = internal.FLD_EGRESS_VRFID
	FldBlockStart        int = internal.FLD_BLOCK_START
	FldBlockEnd          int = internal.FLD_BLOCK_END
	FldBlockStep         int = internal.FLD_BLOCK_STEP
	FldBlockSize         int = internal.FLD_BLOCK_SIZE
	FldClientNwDelayUsec int = internal.FLD_CLIENT_NW_DELAY_USEC
	FldServerNwDelayUsec int = internal.FLD_SERVER_NW_DELAY_USEC
	FldApplLatencyUsec   int = internal.FLD_APPL_LATENCY_USEC
	FldInetFamily        int = internal.FLD_INET_FAMILY
	FldExporterIp        int = internal.FLD_EXPORTER_IP
	FldExporterId        int = internal.FLD_EXPORTER_ID
	FldExporterVersion   int = internal.FLD_EXPORTER_VERSION
	FldSequenceFailures  int = internal.FLD_SEQUENCE_FAILURES
	FldSamplerMode       int = internal.FLD_SAMPLER_MODE
	FldSamplerInterval   int = internal.FLD_SAMPLER_INTERVAL
	FldSamplerId         int = internal.FLD_SAMPLER_ID
	FldCalcDuration      int = internal.FLD_CALC_DURATION
	FldCalcBps           int = internal.FLD_CALC_BPS
	FldCalcPps           int = internal.FLD_CALC_PPS
	FldCalcBpp           int = internal.FLD_CALC_BPP
	FldBrec1             int = internal.FLD_BREC1
	FldPairPort          int = internal.FLD_PAIR_PORT
	FldPairAddr          int = internal.FLD_PAIR_ADDR
	FldPairAddrAlias     int = internal.FLD_PAIR_ADDR_ALIAS
	FldPairAs            int = internal.FLD_PAIR_AS
	FldPairIf            int = internal.FLD_PAIR_IF
	FldPairVlan          int = internal.FLD_PAIR_VLAN
	FldTerm              int = internal.FLD_TERM_
)

var fieldTypes = map[int]any{
	FldFirst:             time.Time{},
	FldLast:              time.Time{},
	FldReceived:          uint64(0),
	FldDoctets:           uint64(0),
	FldDpkts:             uint64(0),
	FldDpktsAlias:        uint64(0),
	FldOutBytes:          uint64(0),
	FldOutPkts:           uint64(0),
	FldOutPktsAlias:      uint64(0),
	FldAggrFlows:         uint64(0),
	FldSrcport:           uint16(0),
	FldDstport:           uint16(0),
	FldTcpFlags:          uint8(0),
	FldTcpFlagsAlias:     uint8(0),
	FldSrcaddr:           net.IP{0},
	FldDstaddr:           net.IP{0},
	FldSrcaddrAlias:      net.IP{0},
	FldDstaddrAlias:      net.IP{0},
	FldIpNextHop:         net.IP{0},
	FldIpNextHopAlias:    net.IP{0},
	FldSrcMask:           uint8(0),
	FldDstMask:           uint8(0),
	FldTos:               uint8(0),
	FldDstTos:            uint8(0),
	FldSrcAS:             uint32(0),
	FldDstAS:             uint32(0),
	FldBgpNextAdjacentAS: uint32(0),
	FldBgpPrevAdjacentAS: uint32(0),
	FldBgpNextHop:        net.IP{0},
	FldProt:              uint8(0),
	FldSrcVlan:           uint16(0),
	FldDstVlan:           uint16(0),
	FldInSrcMac:          MacAddr{0},
	FldOutSrcMac:         MacAddr{0},
	FldInDstMac:          MacAddr{0},
	FldOutDstMac:         MacAddr{0},
	FldMplsLabel:         uint32(0),
	FldInput:             uint32(0),
	FldOutput:            uint32(0),
	FldDir:               uint8(0),
	FldFwdStatus:         uint8(0),
	FldIpRouter:          net.IP{0},
	FldIpRouterAlias:     net.IP{0},
	FldEngineType:        uint8(0),
	FldEngineId:          uint8(0),
	FldEngineTypeAlias:   uint8(0),
	FldEngineIdAlias:     uint8(0),
	FldEventTime:         uint64(0),
	FldConnId:            uint32(0),
	FldIcmpCode:          uint8(0),
	FldIcmpType:          uint8(0),
	FldIcmpCodeAlias:     uint8(0),
	FldIcmpTypeAlias:     uint8(0),
	FldFwXEvent:          uint16(0),
	FldFwEvent:           uint8(0),
	FldXlateSrcIp:        net.IP{0},
	FldXlateDstIp:        net.IP{0},
	FldXlateSrcPort:      uint16(0),
	FldXlateDstPort:      uint16(0),
	FldIngressAclId:      uint32(0),
	FldIngressAceId:      uint32(0),
	FldIngressXaceId:     uint32(0),
	FldIngressAcl:        Acl{},
	FldEgressAclId:       uint32(0),
	FldEgressAceId:       uint32(0),
	FldEgressXaceId:      uint32(0),
	FldEgressAcl:         Acl{},
	FldUsername:          "",
	FldIngressVrfid:      uint32(0),
	FldEventFlag:         uint8(0),
	FldEgressVrfid:       uint32(0),
	FldBlockStart:        uint16(0),
	FldBlockEnd:          uint16(0),
	FldBlockStep:         uint16(0),
	FldBlockSize:         uint16(0),
	FldClientNwDelayUsec: uint64(0),
	FldServerNwDelayUsec: uint64(0),
	FldApplLatencyUsec:   uint64(0),
	FldInetFamily:        uint32(0),
	FldExporterIp:        net.IP{0},
	FldExporterId:        uint32(0),
	FldExporterVersion:   uint32(0),
	FldSequenceFailures:  uint32(0),
	FldSamplerMode:       uint16(0),
	FldSamplerInterval:   uint32(0),
	FldSamplerId:         uint32(0),
	FldCalcDuration:      uint64(0),
	FldCalcBps:           float64(0),
	FldCalcPps:           float64(0),
	FldCalcBpp:           float64(0),
	FldBrec1:             BasicRecord1{},
	FldPairPort:          uint16(0),
	FldPairAddr:          net.IP{0},
	FldPairAddrAlias:     net.IP{0},
	FldPairAs:            uint32(0),
	FldPairIf:            uint16(0),
	FldPairVlan:          uint16(0),
	FldTerm:              uint8(0),
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

func convertToIP(data []byte) net.IP {
	if isAllBytesZero(data[:12]) {
		return net.IP(data[12:]) // IPv4 Address
	}
	return net.IP(data) // IPv6 Address
}

func NewRecord() (Record, error) {
	var r Record
	status := internal.Rec_init(&r.ptr)
	if status == internal.ERR_NOMEM {
		return r, ErrNoMem
	} else if status == internal.ERR_OTHER {
		return r, ErrOther
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
	return ErrOther
}

func callFget(r Record, field int, fieldPtr uintptr) error {
	status := internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
	if status == internal.ERR_NOTSET {
		return ErrNotSet
	}
	return nil
}

func getSimpleDataType(r Record, expectedType any, field int) (any, error) {
	expectedTypeReflect := reflect.TypeOf(expectedType)
	if expectedTypeReflect == nil {
		return nil, ErrUnknownFld
	}
	valPtr := reflect.New(expectedTypeReflect).Interface()
	fieldPtr := unsafe.Pointer(reflect.ValueOf(valPtr).Pointer())
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(valPtr).Elem().Interface(), nil
}

func getIP(r Record, field int) (any, error) {
	ipBuf := make([]byte, 16) // Allocate 16 bytes (IPv6 max size)
	fieldPtr := unsafe.Pointer(&ipBuf[0])
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return convertToIP(ipBuf), nil
}

func getTime(r Record, field int) (any, error) {
	val := int64(0)
	fieldPtr := unsafe.Pointer(&val)
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return time.UnixMilli(val), nil

}

func getMacAddress(r Record, field int) (any, error) {
	val := [6]byte{0}
	fieldPtr := unsafe.Pointer(&val)
	err := callFget(r, field, uintptr(fieldPtr))
	if err != nil {
		return nil, err
	}
	return net.HardwareAddr(val[:]), nil

}

func getBasicRecord1(r Record, field int) (any, error) {
	brec := internal.NewLnf_brec1_t()
	defer internal.DeleteLnf_brec1_t(brec)
	err := callFget(r, field, uintptr(brec.Swigcptr()))

	if err != nil {
		return nil, err
	}

	var output BasicRecord1
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

func getAcl(r Record, field int) (any, error) {
	acl := internal.NewLnf_acl_t()
	defer internal.DeleteLnf_acl_t(acl)
	err := callFget(r, field, uintptr(acl.Swigcptr()))
	if err != nil {
		return nil, err
	}
	return Acl{
		AclId:  acl.GetAcl_id(),
		AceId:  acl.GetAce_id(),
		XaceId: acl.GetXace_id(),
	}, nil
}

func (r Record) GetField(field int) (any, error) {
	expectedType, ok := fieldTypes[field]
	var ret any
	var err error
	if !ok {
		return nil, ErrUnknownFld
	}
	switch expectedType.(type) {
	case uint64, uint32, uint16, uint8, float64:
		ret, err = getSimpleDataType(r, expectedType, field)

	case net.IP:
		ret, err = getIP(r, field)

	case time.Time:
		ret, err = getTime(r, field)

	case MacAddr: // MacAddr
		ret, err = getMacAddress(r, field)

	case BasicRecord1:
		ret, err = getBasicRecord1(r, field)

	case Acl:
		ret, err = getAcl(r, field)

	case string:
		panic("not implemented")
		// Assuming the C function writes to a char buffer
		// buf := make([]byte, 64) // Adjust the size as needed
		// fieldPtr = unsafe.Pointer(&buf[0])
		// internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		// return string(buf), nil

	default:
		ret = nil
		err = ErrUnknownFld
	}
	return ret, err
}
