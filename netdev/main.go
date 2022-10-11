package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"gomodules.xyz/oneliners"
	_ "gomodules.xyz/oneliners"
)

func main() {
	oneliners.FILE()
	la := netlink.NewLinkAttrs()
	oneliners.FILE()
	la.Name = "foo"
	oneliners.FILE()
	mybridge := &netlink.Bridge{LinkAttrs: la}
	oneliners.FILE()
	err := netlink.LinkAdd(mybridge)
	oneliners.FILE(err)
	if err != nil {
		oneliners.FILE()
		fmt.Printf("could not add %s: %v\n", la.Name, err)
	}
	oneliners.FILE()
	eth1, err := netlink.LinkByName("eth0")
	if err != nil {
		oneliners.FILE()
		fmt.Printf(err.Error())
	}
	oneliners.FILE()
	err = netlink.LinkSetMaster(eth1, mybridge)
	oneliners.FILE(err)
}

func main123() {
	lo, _ := netlink.LinkByName("lo")
	addr, _ := netlink.ParseAddr("169.254.169.254/32")
	netlink.AddrAdd(lo, addr)
}
