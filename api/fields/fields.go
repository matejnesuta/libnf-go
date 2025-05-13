package fields

import (
	"libnf-go/internal"
	"net"
	"time"
)

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

type Mpls [10]uint32

const (
	First             int = internal.FLD_FIRST
	Last              int = internal.FLD_LAST
	Received          int = internal.FLD_RECEIVED
	Doctets           int = internal.FLD_DOCTETS
	Dpkts             int = internal.FLD_DPKTS
	DpktsAlias        int = internal.FLD_DPKTS_ALIAS
	OutBytes          int = internal.FLD_OUT_BYTES
	OutPkts           int = internal.FLD_OUT_PKTS
	OutPktsAlias      int = internal.FLD_OUT_PKTS_ALIAS
	AggrFlows         int = internal.FLD_AGGR_FLOWS
	SrcPort           int = internal.FLD_SRCPORT
	DstPort           int = internal.FLD_DSTPORT
	TcpFlags          int = internal.FLD_TCP_FLAGS
	TcpFlagsAlias     int = internal.FLD_TCP_FLAGS_ALIAS
	SrcAddr           int = internal.FLD_SRCADDR
	DstAddr           int = internal.FLD_DSTADDR
	SrcAddrAlias      int = internal.FLD_SRCADDR_ALIAS
	DstAddrAlias      int = internal.FLD_DSTADDR_ALIAS
	IpNextHop         int = internal.FLD_IP_NEXTHOP
	IpNextHopAlias    int = internal.FLD_IP_NEXTHOP_ALIAS
	SrcMask           int = internal.FLD_SRC_MASK
	DstMask           int = internal.FLD_DST_MASK
	Tos               int = internal.FLD_TOS
	DstTos            int = internal.FLD_DST_TOS
	SrcAS             int = internal.FLD_SRCAS
	DstAS             int = internal.FLD_DSTAS
	BgpNextAdjacentAS int = internal.FLD_BGPNEXTADJACENTAS
	BgpPrevAdjacentAS int = internal.FLD_BGPPREVADJACENTAS
	BgpNextHop        int = internal.FLD_BGP_NEXTHOP
	Prot              int = internal.FLD_PROT
	SrcVlan           int = internal.FLD_SRC_VLAN
	DstVlan           int = internal.FLD_DST_VLAN
	InSrcMac          int = internal.FLD_IN_SRC_MAC
	OutSrcMac         int = internal.FLD_OUT_SRC_MAC
	InDstMac          int = internal.FLD_IN_DST_MAC
	OutDstMac         int = internal.FLD_OUT_DST_MAC
	MplsLabel         int = internal.FLD_MPLS_LABEL
	Input             int = internal.FLD_INPUT
	Output            int = internal.FLD_OUTPUT
	Dir               int = internal.FLD_DIR
	FwdStatus         int = internal.FLD_FWD_STATUS
	IpRouter          int = internal.FLD_IP_ROUTER
	IpRouterAlias     int = internal.FLD_IP_ROUTER_ALIAS
	EngineType        int = internal.FLD_ENGINE_TYPE
	EngineId          int = internal.FLD_ENGINE_ID
	EngineTypeAlias   int = internal.FLD_ENGINE_TYPE_ALIAS
	EngineIdAlias     int = internal.FLD_ENGINE_ID_ALIAS
	EventTime         int = internal.FLD_EVENT_TIME
	ConnId            int = internal.FLD_CONN_ID
	IcmpCode          int = internal.FLD_ICMP_CODE
	IcmpType          int = internal.FLD_ICMP_TYPE
	IcmpCodeAlias     int = internal.FLD_ICMP_CODE_ALIAS
	IcmpTypeAlias     int = internal.FLD_ICMP_TYPE_ALIAS
	FwXEvent          int = internal.FLD_FW_XEVENT
	FwEvent           int = internal.FLD_FW_EVENT
	XlateSrcIp        int = internal.FLD_XLATE_SRC_IP
	XlateDstIp        int = internal.FLD_XLATE_DST_IP
	XlateSrcPort      int = internal.FLD_XLATE_SRC_PORT
	XlateDstPort      int = internal.FLD_XLATE_DST_PORT
	IngressAclId      int = internal.FLD_INGRESS_ACL_ID
	IngressAceId      int = internal.FLD_INGRESS_ACE_ID
	IngressXaceId     int = internal.FLD_INGRESS_XACE_ID
	IngressAcl        int = internal.FLD_INGRESS_ACL
	EgressAclId       int = internal.FLD_EGRESS_ACL_ID
	EgressAceId       int = internal.FLD_EGRESS_ACE_ID
	EgressXaceId      int = internal.FLD_EGRESS_XACE_ID
	EgressAcl         int = internal.FLD_EGRESS_ACL
	Username          int = internal.FLD_USERNAME
	IngressVrfid      int = internal.FLD_INGRESS_VRFID
	EventFlag         int = internal.FLD_EVENT_FLAG
	EgressVrfid       int = internal.FLD_EGRESS_VRFID
	BlockStart        int = internal.FLD_BLOCK_START
	BlockEnd          int = internal.FLD_BLOCK_END
	BlockStep         int = internal.FLD_BLOCK_STEP
	BlockSize         int = internal.FLD_BLOCK_SIZE
	ClientNwDelayUsec int = internal.FLD_CLIENT_NW_DELAY_USEC
	ServerNwDelayUsec int = internal.FLD_SERVER_NW_DELAY_USEC
	ApplLatencyUsec   int = internal.FLD_APPL_LATENCY_USEC
	InetFamily        int = internal.FLD_INET_FAMILY
	ExporterIp        int = internal.FLD_EXPORTER_IP
	ExporterId        int = internal.FLD_EXPORTER_ID
	ExporterVersion   int = internal.FLD_EXPORTER_VERSION
	SequenceFailures  int = internal.FLD_SEQUENCE_FAILURES
	SamplerMode       int = internal.FLD_SAMPLER_MODE
	SamplerInterval   int = internal.FLD_SAMPLER_INTERVAL
	SamplerId         int = internal.FLD_SAMPLER_ID
	CalcDuration      int = internal.FLD_CALC_DURATION
	CalcBps           int = internal.FLD_CALC_BPS
	CalcPps           int = internal.FLD_CALC_PPS
	CalcBpp           int = internal.FLD_CALC_BPP
	Brec1             int = internal.FLD_BREC1
	PairPort          int = internal.FLD_PAIR_PORT
	PairAddr          int = internal.FLD_PAIR_ADDR
	PairAddrAlias     int = internal.FLD_PAIR_ADDR_ALIAS
	PairAs            int = internal.FLD_PAIR_AS
	PairIf            int = internal.FLD_PAIR_IF
	PairVlan          int = internal.FLD_PAIR_VLAN
	Term              int = internal.FLD_TERM_
)

var FieldTypes = map[int]any{
	First:             time.Time{},
	Last:              time.Time{},
	Received:          uint64(0),
	Doctets:           uint64(0),
	Dpkts:             uint64(0),
	DpktsAlias:        uint64(0),
	OutBytes:          uint64(0),
	OutPkts:           uint64(0),
	OutPktsAlias:      uint64(0),
	AggrFlows:         uint64(0),
	SrcPort:           uint16(0),
	DstPort:           uint16(0),
	TcpFlags:          uint8(0),
	TcpFlagsAlias:     uint8(0),
	SrcAddr:           net.IP{0},
	DstAddr:           net.IP{0},
	SrcAddrAlias:      net.IP{0},
	DstAddrAlias:      net.IP{0},
	IpNextHop:         net.IP{0},
	IpNextHopAlias:    net.IP{0},
	SrcMask:           uint8(0),
	DstMask:           uint8(0),
	Tos:               uint8(0),
	DstTos:            uint8(0),
	SrcAS:             uint32(0),
	DstAS:             uint32(0),
	BgpNextAdjacentAS: uint32(0),
	BgpPrevAdjacentAS: uint32(0),
	BgpNextHop:        net.IP{0},
	Prot:              uint8(0),
	SrcVlan:           uint16(0),
	DstVlan:           uint16(0),
	InSrcMac:          net.HardwareAddr{0},
	OutSrcMac:         net.HardwareAddr{0},
	InDstMac:          net.HardwareAddr{0},
	OutDstMac:         net.HardwareAddr{0},
	MplsLabel:         Mpls{},
	Input:             uint32(0),
	Output:            uint32(0),
	Dir:               uint8(0),
	FwdStatus:         uint8(0),
	IpRouter:          net.IP{0},
	IpRouterAlias:     net.IP{0},
	EngineType:        uint8(0),
	EngineId:          uint8(0),
	EngineTypeAlias:   uint8(0),
	EngineIdAlias:     uint8(0),
	EventTime:         uint64(0),
	ConnId:            uint32(0),
	IcmpCode:          uint8(0),
	IcmpType:          uint8(0),
	IcmpCodeAlias:     uint8(0),
	IcmpTypeAlias:     uint8(0),
	FwXEvent:          uint16(0),
	FwEvent:           uint8(0),
	XlateSrcIp:        net.IP{0},
	XlateDstIp:        net.IP{0},
	XlateSrcPort:      uint16(0),
	XlateDstPort:      uint16(0),
	IngressAclId:      uint32(0),
	IngressAceId:      uint32(0),
	IngressXaceId:     uint32(0),
	IngressAcl:        Acl{},
	EgressAclId:       uint32(0),
	EgressAceId:       uint32(0),
	EgressXaceId:      uint32(0),
	EgressAcl:         Acl{},
	Username:          "",
	IngressVrfid:      uint32(0),
	EventFlag:         uint8(0),
	EgressVrfid:       uint32(0),
	BlockStart:        uint16(0),
	BlockEnd:          uint16(0),
	BlockStep:         uint16(0),
	BlockSize:         uint16(0),
	ClientNwDelayUsec: uint64(0),
	ServerNwDelayUsec: uint64(0),
	ApplLatencyUsec:   uint64(0),
	InetFamily:        uint32(0),
	ExporterIp:        net.IP{0},
	ExporterId:        uint32(0),
	ExporterVersion:   uint32(0),
	SequenceFailures:  uint32(0),
	SamplerMode:       uint16(0),
	SamplerInterval:   uint32(0),
	SamplerId:         uint32(0),
	CalcDuration:      uint64(0),
	CalcBps:           float64(0),
	CalcPps:           float64(0),
	CalcBpp:           float64(0),
	Brec1:             BasicRecord1{},
	PairPort:          uint16(0),
	PairAddr:          net.IP{0},
	PairAddrAlias:     net.IP{0},
	PairAs:            uint32(0),
	PairIf:            uint16(0),
	PairVlan:          uint16(0),
	Term:              uint8(0),
}

type FldDataType interface {
	uint8 | uint16 | uint32 | uint64 | net.IP | time.Time | net.HardwareAddr | BasicRecord1 | Acl | Mpls | string
}
