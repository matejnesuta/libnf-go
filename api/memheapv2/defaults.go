package memheapv2

import (
	"libnf-go/api/fields"
)

const (
	SortNone int = 0
	SortAsc  int = 16
	SortDesc int = 32
)

const (
	AggrAuto int = 0
	AggrMin  int = 1
	AggrMax  int = 2
	AggrSum  int = 3
	AggrOr   int = 4
	AggrKey  int = 8
)

var defaults = map[int][2]int{
	fields.First:             {AggrMin, SortAsc},
	fields.Last:              {AggrMax, SortDesc},
	fields.Received:          {AggrMax, SortAsc},
	fields.Doctets:           {AggrSum, SortDesc},
	fields.Dpkts:             {AggrSum, SortDesc},
	fields.DpktsAlias:        {AggrSum, SortDesc},
	fields.OutBytes:          {AggrSum, SortDesc},
	fields.OutPkts:           {AggrSum, SortDesc},
	fields.OutPktsAlias:      {AggrSum, SortDesc},
	fields.AggrFlows:         {AggrSum, SortDesc},
	fields.SrcPort:           {AggrKey, SortAsc},
	fields.DstPort:           {AggrKey, SortAsc},
	fields.TcpFlags:          {AggrOr, SortAsc},
	fields.TcpFlagsAlias:     {AggrOr, SortAsc},
	fields.SrcAddr:           {AggrKey, SortAsc},
	fields.DstAddr:           {AggrKey, SortAsc},
	fields.SrcAddrAlias:      {AggrKey, SortAsc},
	fields.DstAddrAlias:      {AggrKey, SortAsc},
	fields.IpNextHop:         {AggrKey, SortAsc},
	fields.IpNextHopAlias:    {AggrKey, SortAsc},
	fields.SrcMask:           {AggrKey, SortAsc},
	fields.DstMask:           {AggrKey, SortAsc},
	fields.Tos:               {AggrKey, SortAsc},
	fields.DstTos:            {AggrKey, SortAsc},
	fields.SrcAS:             {AggrKey, SortAsc},
	fields.DstAS:             {AggrKey, SortAsc},
	fields.BgpNextAdjacentAS: {AggrKey, SortAsc},
	fields.BgpPrevAdjacentAS: {AggrKey, SortAsc},
	fields.BgpNextHop:        {AggrKey, SortAsc},
	fields.Prot:              {AggrKey, SortAsc},
	fields.SrcVlan:           {AggrKey, SortAsc},
	fields.DstVlan:           {AggrKey, SortAsc},
	fields.InSrcMac:          {AggrKey, SortAsc},
	fields.OutSrcMac:         {AggrKey, SortAsc},
	fields.InDstMac:          {AggrKey, SortAsc},
	fields.OutDstMac:         {AggrKey, SortAsc},
	fields.MplsLabel:         {AggrKey, SortNone},
	fields.Input:             {AggrKey, SortAsc},
	fields.Output:            {AggrKey, SortAsc},
	fields.Dir:               {AggrKey, SortAsc},
	fields.FwdStatus:         {AggrKey, SortAsc},
	fields.IpRouter:          {AggrKey, SortAsc},
	fields.IpRouterAlias:     {AggrKey, SortAsc},
	fields.EngineType:        {AggrKey, SortAsc},
	fields.EngineId:          {AggrKey, SortAsc},
	fields.EngineTypeAlias:   {AggrKey, SortAsc},
	fields.EngineIdAlias:     {AggrKey, SortAsc},
	fields.EventTime:         {AggrMin, SortAsc},
	fields.ConnId:            {AggrKey, SortAsc},
	fields.IcmpCode:          {AggrKey, SortAsc},
	fields.IcmpType:          {AggrKey, SortAsc},
	fields.IcmpCodeAlias:     {AggrKey, SortAsc},
	fields.IcmpTypeAlias:     {AggrKey, SortAsc},
	fields.FwXEvent:          {AggrKey, SortAsc},
	fields.FwEvent:           {AggrKey, SortAsc},
	fields.XlateSrcIp:        {AggrKey, SortAsc},
	fields.XlateDstIp:        {AggrKey, SortAsc},
	fields.XlateSrcPort:      {AggrKey, SortAsc},
	fields.XlateDstPort:      {AggrKey, SortAsc},
	fields.IngressAclId:      {AggrKey, SortAsc},
	fields.IngressAceId:      {AggrKey, SortAsc},
	fields.IngressXaceId:     {AggrKey, SortAsc},
	fields.IngressAcl:        {AggrKey, SortAsc},
	fields.EgressAceId:       {AggrKey, SortAsc},
	fields.EgressXaceId:      {AggrKey, SortAsc},
	fields.EgressAclId:       {AggrKey, SortAsc},
	fields.EgressAcl:         {AggrKey, SortAsc},
	fields.Username:          {AggrKey, SortAsc},
	fields.IngressVrfid:      {AggrKey, SortAsc},
	fields.EventFlag:         {AggrKey, SortAsc},
	fields.EgressVrfid:       {AggrKey, SortAsc},
	fields.BlockStart:        {AggrKey, SortAsc},
	fields.BlockEnd:          {AggrKey, SortAsc},
	fields.BlockStep:         {AggrKey, SortAsc},
	fields.BlockSize:         {AggrKey, SortAsc},
	fields.ClientNwDelayUsec: {AggrKey, SortAsc},
	fields.ServerNwDelayUsec: {AggrKey, SortAsc},
	fields.ApplLatencyUsec:   {AggrKey, SortAsc},
	fields.InetFamily:        {AggrKey, SortAsc},
	fields.ExporterIp:        {AggrKey, SortAsc},
	fields.ExporterId:        {AggrKey, SortAsc},
	fields.ExporterVersion:   {AggrKey, SortAsc},
	fields.SequenceFailures:  {AggrSum, SortAsc},
	fields.SamplerMode:       {AggrKey, SortAsc},
	fields.SamplerInterval:   {AggrKey, SortAsc},
	fields.SamplerId:         {AggrKey, SortAsc},
	fields.CalcDuration:      {AggrSum, SortAsc},
	fields.CalcBps:           {AggrSum, SortDesc},
	fields.CalcPps:           {AggrSum, SortDesc},
	fields.CalcBpp:           {AggrSum, SortAsc},
	fields.Brec1:             {AggrKey, SortNone},
	fields.PairPort:          {AggrKey, SortNone},
	fields.PairAddr:          {AggrKey, SortNone},
	fields.PairAddrAlias:     {AggrKey, SortNone},
	fields.PairAs:            {AggrKey, SortNone},
	fields.PairIf:            {AggrKey, SortNone},
	fields.PairVlan:          {AggrKey, SortNone},
}

var dependencies = map[int][]int{
	fields.CalcDuration: {fields.First, fields.Last},
	fields.CalcBps:      {fields.Doctets, fields.CalcDuration},
	fields.CalcPps:      {fields.Dpkts, fields.CalcDuration},
	fields.CalcBpp:      {fields.Dpkts, fields.Doctets},
}

var pairFields = map[int][2]int{
	fields.PairPort:      {fields.SrcPort, fields.DstPort},
	fields.PairAddr:      {fields.SrcAddr, fields.DstAddr},
	fields.PairAddrAlias: {fields.SrcAddrAlias, fields.DstAddrAlias},
	fields.PairAs:        {fields.SrcAS, fields.DstAS},
	fields.PairIf:        {fields.Input, fields.Output},
	fields.PairVlan:      {fields.SrcVlan, fields.DstVlan},
}
