package main

import (
	"fmt"
	libnf "libnf/api"
)

func main() {
	var ptr libnf.File
	err := ptr.OpenRead("testfiles/nfcapd.201705281555", false, false)
	if err != nil {
		fmt.Println(err)
	}
	nfdump_version, err := ptr.GetNfdumpVersion()
	ident, err := ptr.GetIdent()
	bytes, err := ptr.GetBytes()
	packets, err := ptr.GetPackets()
	flows, err := ptr.GetFlows()
	fmt.Println(nfdump_version)
	fmt.Println(ident)
	fmt.Println(bytes)
	fmt.Println(packets)
	fmt.Println(flows)
	ptr.Close()

}
