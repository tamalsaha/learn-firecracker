package main

import (
	"fmt"
	"net"
)

func main() {
	instanceID := 0
	VMS_NETWORK_PREFIX := "172.26.0"

	// binary.Write(a, binary.LittleEndian, myInt)
	ip_1 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, (instanceID+1)*2)
	ip_2 := fmt.Sprintf("%s.%d", VMS_NETWORK_PREFIX, (instanceID+1)*2+1)

	fmt.Println("instanceID:", instanceID)
	fmt.Println("IP_1:", ip_1, MacAddr(net.ParseIP(ip_1).To4()))
	fmt.Println("IP_2:", ip_2, MacAddr(net.ParseIP(ip_2).To4()))

}

const hexDigit = "0123456789abcdef"

func MacAddr(b []byte) string {
	s := make([]byte, len(b)*3)
	for i, tn := range b {
		s[i*3], s[i*3+1], s[i*3+2] = ':', hexDigit[tn>>4], hexDigit[tn&0xf]
	}
	return "AA:FF" + string(s)
}
