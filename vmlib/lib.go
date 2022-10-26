package vmlib

import (
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
