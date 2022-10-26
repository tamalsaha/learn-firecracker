package main

import (
	"fmt"
	"gomodules.xyz/go-sh"
)

// EGRESS_IFACE=`ip route get 1.1.1.1 |grep uid |sed 's/.* dev \([^ ]*\) .*/\1/'`	#
func main() {
	egressIface, err := sh.
		Command("ip", "route", "get", "1.1.1.1").
		Command("grep", "uid").
		Command("sed", `s/.* dev \([^ ]*\) .*/\1/`).
		Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(egressIface))
}

func GetEgressInterface() (string, error) {
	egressIface, err := sh.
		Command("ip", "route", "get", "1.1.1.1").
		Command("grep", "uid").
		Command("sed", `s/.* dev \([^ ]*\) .*/\1/`).
		Output()
	return string(egressIface), err
}
