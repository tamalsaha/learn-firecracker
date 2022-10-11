package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func main() {
	la := netlink.NewLinkAttrs()
	la.Name = "foo"
	mybridge := &netlink.Bridge{LinkAttrs: la}
	err := netlink.LinkAdd(mybridge)
	if err != nil {
		fmt.Printf("could not add %s: %v\n", la.Name, err)
	}
	eth1, _ := netlink.LinkByName("bond0")
	netlink.LinkSetMaster(eth1, mybridge)
}

func main123() {
	lo, _ := netlink.LinkByName("lo")
	addr, _ := netlink.ParseAddr("169.254.169.254/32")
	netlink.AddrAdd(lo, addr)
}
