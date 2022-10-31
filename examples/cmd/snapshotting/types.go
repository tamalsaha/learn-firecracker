package main

type NetworkConfig struct {
	Version   int                       `json:"version"`
	Ethernets map[string]EthernetConfig `json:"ethernets"`
}

type EthernetConfig struct {
	Match       EthernetMatcher `json:"match"`
	Addresses   []string        `json:"addresses"`
	Gateway4    string          `json:"gateway4,omitempty"`
	Nameservers *Nameservers    `json:"nameservers,omitempty"`
}

type EthernetMatcher struct {
	Macaddress string `json:"macaddress"`
}

type Nameservers struct {
	Addresses []string `json:"addresses"`
}

////

type Metadata struct {
	InstanceID    string `json:"instance-id"`
	LocalHostname string `json:"local-hostname"`
}

////

type UserData struct {
	Users []User `json:"users"`
}

type User struct {
	Name              string   `json:"name"`
	Gecos             string   `json:"gecos,omitempty"`
	Shell             string   `json:"shell,omitempty"`
	Groups            string   `json:"groups,omitempty"`
	Sudo              string   `json:"sudo,omitempty"`
	SSHAuthorizedKeys []string `json:"ssh_authorized_keys,omitempty"`
}

///

type MMDSConfig struct {
	Latest LatestConfig `json:"latest"`
}

type LatestConfig struct {
	MetaData interface{} `json:"meta-data,omitempty"`
	// UserData string `json:"user-data,omitempty"`
}
