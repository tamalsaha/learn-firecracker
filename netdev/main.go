package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"gomodules.xyz/oneliners"
	_ "gomodules.xyz/oneliners"
)

/*
sudo ip tuntap add tap0 mode tap
sudo ip addr add 172.16.0.1/24 dev tap0
sudo ip link set tap0 up

*/

func main() {
	nl, err := CreateTap("tap0", 0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}

}

func CreateTap(name string, mtu int, ownerUID, ownerGID int) (netlink.Link, error) {
	tapLinkAttrs := netlink.NewLinkAttrs()
	tapLinkAttrs.Name = name
	tapLink := &netlink.Tuntap{
		LinkAttrs: tapLinkAttrs,

		// We want a tap device (L2) as opposed to a tun (L3)
		Mode: netlink.TUNTAP_MODE_TAP,

		// Firecracker does not support multiqueue tap devices at this time:
		// https://github.com/firecracker-microvm/firecracker/issues/750
		Queues: 1,

		Flags: netlink.TUNTAP_ONE_QUEUE | // single queue tap device
			netlink.TUNTAP_VNET_HDR, // parse vnet headers added by the vm's virtio_net implementation
	}

	err := netlink.LinkAdd(tapLink)
	if err != nil {
		return nil, fmt.Errorf("failed to create tap device: %w", err)
	}

	//for _, tapFd := range tapLink.Fds {
	//	err = unix.IoctlSetInt(int(tapFd.Fd()), unix.TUNSETOWNER, ownerUID)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to set tap %s owner to uid %d: %w",
	//			name, ownerUID, err)
	//
	//	}
	//
	//	err = unix.IoctlSetInt(int(tapFd.Fd()), unix.TUNSETGROUP, ownerGID)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to set tap %s group to gid %d: %w",
	//			name, ownerGID, err)
	//
	//	}
	//}

	//err = netlink.LinkSetMTU(tapLink, mtu)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to set tap device MTU to %d: %w", mtu, err)
	//}

	addr, err := netlink.ParseAddr("172.25.0.1/24")
	if err != nil {
		return nil, err
	}
	err = netlink.AddrAdd(tapLink, addr)
	if err != nil {
		return nil, err
	}

	err = netlink.LinkSetUp(tapLink)
	if err != nil {
		return nil, fmt.Errorf("failed to set tap up: %w", err)
	}

	return tapLink, nil
}

func main__() {
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
	err = netlink.LinkSetMaster(mybridge, eth1)
	oneliners.FILE(err)
}

func main123() {
	lo, _ := netlink.LinkByName("lo")
	addr, _ := netlink.ParseAddr("169.254.169.254/32")
	netlink.AddrAdd(lo, addr)
}
