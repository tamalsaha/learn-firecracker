package main

import (
	"fmt"
	"gomodules.xyz/go-sh"
	"net"
)

// EGRESS_IFACE=`ip route get 1.1.1.1 |grep uid |sed 's/.* dev \([^ ]*\) .*/\1/'`	#
func main__() {
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

func main() {
	msk := net.CIDRMask(24, 8*net.IPv4len)
	fmt.Println(msk.String())
}
