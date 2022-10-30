package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"gomodules.xyz/go-sh"
)

func GetEgressInterface() (string, error) {
	egressIface, err := sh.
		Command("ip", "route", "get", "1.1.1.1").
		Command("grep", "uid").
		Command("sed", `s/.* dev \([^ ]*\) .*/\1/`).
		Output()
	return strings.TrimSpace(string(egressIface)), err
}

// Host iface bond0
// detect using
// EGRESS_IFACE=`ip route get 8.8.8.8 |grep uid |sed 's/.* dev \([^ ]*\) .*/\1/'`
func SetupIPTables(iface, tapDev string) error {
	/*
		sudo sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
		sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
		sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
		sudo iptables -A FORWARD -i tap0 -o eth0 -j ACCEPT
	*/

	// sudo sysctl -w net.ipv4.ip_forward=1 > /dev/null
	filename := "/proc/sys/net/ipv4/ip_forward"
	if err := os.WriteFile(filename, []byte("1"), 0o644); err != nil {
		return errors.Wrapf(err, "failed to write file: %s", filename)
	}

	tbl, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv6), iptables.Timeout(5))
	if err != nil {
		return errors.Wrap(err, "failed to construct iptables obj")
	}
	err = tbl.AppendUnique("nat", "POSTROUTING", "-o", iface, "-j", "MASQUERADE")
	if err != nil {
		return errors.Wrap(err, "failed to exec: "+"iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE")
	}

	/*
		sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
		sudo iptables -A FORWARD -i tap0 -o eth0 -j ACCEPT
	*/
	err = tbl.AppendUnique("filter", "FORWARD", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT")
	if err != nil {
		return errors.Wrap(err, "failed to exec: "+"iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT")
	}
	err = tbl.AppendUnique("filter", "FORWARD", "-i", tapDev, "-o", iface, "-j", "ACCEPT")
	if err != nil {
		return errors.Wrapf(err, "failed to exec: "+"-i %s -o %s -j ACCEPT", tapDev, iface)
	}

	return nil
}

/*
sudo ip tuntap add tap0 mode tap
sudo ip addr add 172.16.0.1/24 dev tap0
sudo ip link set tap0 up
*/
func CreateTap(name, cidr string) (netlink.Link, error) {
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

	if cidr != "" {
		addr, err := netlink.ParseAddr(cidr) // "172.20.0.1/24"
		if err != nil {
			return nil, err
		}
		err = netlink.AddrAdd(tapLink, addr)
		if err != nil && !errors.Is(err, syscall.EEXIST) {
			// oneliners.FILE(err)
			// fmt.Println("errors.Is(err, syscall.EEXIST) = ", errors.Is(err, syscall.EEXIST))
			// fmt.Printf("%T\n", err)
			return nil, err
		}
	}

	err = netlink.LinkSetUp(tapLink)
	if err != nil {
		return nil, fmt.Errorf("failed to set tap up: %w", err)
	}

	return tapLink, nil
}

const hexDigit = "0123456789abcdef"

func MacAddr(b []byte) string {
	s := make([]byte, len(b)*3)
	for i, tn := range b {
		s[i*3], s[i*3+1], s[i*3+2] = ':', hexDigit[tn>>4], hexDigit[tn&0xf]
	}
	return "AA:FC" + string(s)
}

func main_123() {
	instanceID := 0
	VMS_NETWORK_PREFIX := "172.26.0"

	egressIface, err := GetEgressInterface()
	if err != nil {
		panic(err)
	}

	// binary.Write(a, binary.LittleEndian, myInt)
	ip0 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, (instanceID+1)*2)
	ip1 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, (instanceID+1)*2+1)

	fmt.Println("instanceID:", instanceID)
	eth0Mac := MacAddr(net.ParseIP(ip0).To4())
	fmt.Println("ip0:", ip0, eth0Mac)
	eth1Mac := MacAddr(net.ParseIP(ip1).To4())
	fmt.Println("ip1:", ip1, eth1Mac)

	tap0 := fmt.Sprintf("tap%d", (instanceID+1)*2)
	tap1 := fmt.Sprintf("tap%d", (instanceID+1)*2+1)

	fmt.Println(tap0, tap1)
	if _, err := CreateTap(tap0, ""); err != nil {
		panic(err)
	}
	if _, err := CreateTap(tap1, fmt.Sprintf("%s/%d", ip0, VMS_NETWORK_SUBNET)); err != nil {
		panic(err)
	}
	if err = SetupIPTables(egressIface, tap1); err != nil {
		panic(err)
	}

	_, err = BuildNetCfg(eth0Mac, eth1Mac, ip0, ip1)
	if err != nil {
		panic(err)
	}
}
