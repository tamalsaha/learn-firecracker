package main

import (
	"fmt"
	"gomodules.xyz/go-sh"
)

// EGRESS_IFACE=`ip route get 8.8.8.8 |grep uid |sed 's/.* dev \([^ ]*\) .*/\1/'`	#
func main() {
	egressIface, err := sh.
		Command("ip", "route", "get", "8.8.8.8").
		Command("grep", "uid").
		Command("sed", `'s/.* dev \([^ ]*\) .*/\1/'`).
		Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(egressIface))
}
