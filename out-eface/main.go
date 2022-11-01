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

func main____() {
	msk := net.CIDRMask(24, 8*net.IPv4len)
	fmt.Println(msk.String())
}

const VMS_NETWORK_PREFIX = "172.26.0"

const hexDigit = "0123456789abcdef"

func MacAddr(b []byte) string {
	s := make([]byte, len(b)*3)
	for i, tn := range b {
		s[i*3], s[i*3+1], s[i*3+2] = ':', hexDigit[tn>>4], hexDigit[tn&0xf]
	}
	return "AA:FC" + string(s)
}

func main() {
	instanceID := 1

	//egressIface, err := GetEgressInterface()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(egressIface)

	// binary.Write(a, binary.LittleEndian, myInt)
	//ip0 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, (instanceID+1)*2)
	//ip1 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, (instanceID+1)*2+1)

	ip0 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, instanceID*4+1)
	ip1 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, instanceID*4+2)

	fmt.Println("instanceID:", instanceID)
	eth0Mac := MacAddr(net.ParseIP(ip0).To4())
	fmt.Println("ip0:", ip0, eth0Mac)
	eth1Mac := MacAddr(net.ParseIP(ip1).To4())
	fmt.Println("ip1:", ip1, eth1Mac)

	//tap0 := fmt.Sprintf("tap%d", (instanceID+1)*2)
	//tap1 := fmt.Sprintf("tap%d", (instanceID+1)*2+1)
	tap0 := fmt.Sprintf("tap%d", instanceID*4+1)
	tap1 := fmt.Sprintf("tap%d", instanceID*4+2)

	fmt.Println(tap0, tap1)
	//if _, err := CreateTap(tap0, ""); err != nil {
	//	panic(err)
	//}
	//if _, err := CreateTap(tap1, fmt.Sprintf("%s/%d", ip0, VMS_NETWORK_SUBNET)); err != nil {
	//	panic(err)
	//}
	//if err = SetupIPTables(egressIface, tap1); err != nil {
	//	panic(err)
	//}
	//
	//_, err = BuildNetCfg(eth0Mac, eth1Mac, ip0, ip1)
	//if err != nil {
	//	panic(err)
	//}
}
