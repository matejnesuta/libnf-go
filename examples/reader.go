package examples

import (
	"fmt"

	"github.com/matejnesuta/libnf-go/api/fields"
	"github.com/matejnesuta/libnf-go/api/file"
	"github.com/matejnesuta/libnf-go/api/record"
)

func Reader() {
	var ptr file.File
	err := ptr.OpenRead("api/testfiles/profiling.tmp", false, false)

	if err != nil {
		fmt.Println(err)
	}
	defer ptr.Close()

	rec, err := record.NewRecord()
	if err != nil {
		fmt.Println(err)
	}
	defer rec.Free()

	for {
		err = ptr.GetNextRecord(&rec)
		if err != nil {
			break
		}
		val, _ := rec.GetField(fields.Brec1)
		brec, ok := val.(fields.BasicRecord1)
		if !ok {
			panic("Error: Not a BasicRecord1")
		}

		val, _ = rec.GetField(fields.TcpFlags)
		flags, ok := val.(uint8)
		if !ok {
			panic("Error: Not a TcpFlags")
		}

		val, _ = rec.GetField(fields.CalcBpp)
		bpp, ok := val.(float64)
		if !ok {
			panic("Error: Not a CalcBpp")
		}
		val, _ = rec.GetField(fields.CalcBps)
		bps, ok := val.(float64)
		if !ok {
			panic("Error: Not a CalcBps")
		}
		val, _ = rec.GetField(fields.CalcPps)
		pps, ok := val.(float64)
		if !ok {
			panic("Error: Not a CalcPps")
		}

		fmt.Print(brec.First.Format("2006-01-02 15:04:05"), " ")
		fmt.Printf("| %.3f | %3d | %-15s:%-5d | %-15s:%-5d | %6d | %8d | %4d | %d | %f | %f | %f\n", brec.Last.Sub(brec.First).Seconds(), brec.Prot, brec.SrcAddr, brec.SrcPort, brec.DstAddr, brec.DstPort, brec.Pkts, brec.Bytes, brec.Flows, flags, bpp, bps, pps)
		// printBrec(&brec)
	}
}
