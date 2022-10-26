package main

import (
	"os"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"gomodules.xyz/go-sh"
)

func GetEgressInterface() (string, error) {
	egressIface, err := sh.
		Command("ip", "route", "get", "1.1.1.1").
		Command("grep", "uid").
		Command("sed", `s/.* dev \([^ ]*\) .*/\1/`).
		Output()
	return string(egressIface), err
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
	if err := os.WriteFile(filename, []byte("1"), 0644); err != nil {
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
