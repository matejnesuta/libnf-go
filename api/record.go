package libnf

import (
	"libnf/internal"
	"net"
	"time"
	"unsafe"
)

type Record struct {
	ptr       uintptr
	allocated bool
}

type IPAddr net.IP
type MacAddr net.HardwareAddr
type BasicRecord1 uintptr
type Timestamp time.Time

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
	FldFirst:             Timestamp{},
	FldLast:              Timestamp{},
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
	FldSrcaddr:           IPAddr{0},
	FldDstaddr:           IPAddr{0},
	FldSrcaddrAlias:      IPAddr{0},
	FldDstaddrAlias:      IPAddr{0},
	FldIpNextHop:         IPAddr{0},
	FldIpNextHopAlias:    IPAddr{0},
	FldSrcMask:           uint8(0),
	FldDstMask:           uint8(0),
	FldTos:               uint8(0),
	FldDstTos:            uint8(0),
	FldSrcAS:             uint32(0),
	FldDstAS:             uint32(0),
	FldBgpNextAdjacentAS: uint32(0),
	FldBgpPrevAdjacentAS: uint32(0),
	FldBgpNextHop:        IPAddr{0},
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
	FldIpRouter:          IPAddr{0},
	FldIpRouterAlias:     IPAddr{0},
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
	FldXlateSrcIp:        IPAddr{0},
	FldXlateDstIp:        IPAddr{0},
	FldXlateSrcPort:      uint16(0),
	FldXlateDstPort:      uint16(0),
	FldIngressAclId:      uint32(0),
	FldIngressAceId:      uint32(0),
	FldIngressXaceId:     uint32(0),
	FldIngressAcl:        uint32(0),
	FldEgressAclId:       uint32(0),
	FldEgressAceId:       uint32(0),
	FldEgressXaceId:      uint32(0),
	FldEgressAcl:         uint32(0),
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
	FldExporterIp:        IPAddr{0},
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
	FldBrec1:             BasicRecord1(0),
	FldPairPort:          uint16(0),
	FldPairAddr:          IPAddr{0},
	FldPairAddrAlias:     IPAddr{0},
	FldPairAs:            uint32(0),
	FldPairIf:            uint16(0),
	FldPairVlan:          uint16(0),
	FldTerm:              uint8(0),
}

func IsAllBytesZero(data []byte) bool {
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

func (r Record) GetField(field int) (any, error) {
	expectedType, ok := fieldTypes[field]
	if !ok {
		return nil, ErrUnknownFld
	}
	var fieldPtr unsafe.Pointer
	switch expectedType.(type) {
	case uint64:
		val := uint64(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return val, nil

	case uint32:
		val := uint32(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return val, nil

	case uint16:
		val := uint16(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return val, nil

	case uint8:
		val := uint8(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return val, nil

	case int64:
		val := int64(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return val, nil

	case float64:
		val := float64(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return val, nil

	case IPAddr:
		// Assume we don't know if it's IPv4 or IPv6
		ipBuf := make([]byte, 16) // Allocate 16 bytes (IPv6 max size)
		fieldPtr = unsafe.Pointer(&ipBuf[0])

		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		if IsAllBytesZero(ipBuf[:12]) {
			return net.IP(ipBuf[12:]), nil // IPv4 Address
		}
		return net.IP(ipBuf), nil // IPv6 Address

	case Timestamp:
		val := uint64(0)
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return time.UnixMilli(int64(val)), nil

	case MacAddr: // MacAddr
		val := [6]byte{0}
		fieldPtr = unsafe.Pointer(&val)
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return net.HardwareAddr(val[:]), nil

	case string:
		// Assuming the C function writes to a char buffer
		buf := make([]byte, 64) // Adjust the size as needed
		fieldPtr = unsafe.Pointer(&buf[0])
		internal.Rec_fget(r.ptr, field, uintptr(fieldPtr))
		return string(buf), nil

	default:
		return nil, ErrUnknownFld
	}
}
