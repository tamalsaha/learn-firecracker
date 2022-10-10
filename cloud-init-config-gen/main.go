package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sigs.k8s.io/yaml"
	"strings"
)

func main__() {
	keys, err := getSSHPubKeys("tamalsaha")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(keys))
}

func main() {
	nc, err := BuildNetCfg()
	if err != nil {
		panic(err)
	}
	fmt.Println(nc)

	//err := buildData("i123")
	//if err != nil {
	//	panic(err)
	//}
}

func BuildNetCfg() (string, error) {
	/*
		version: 2
		ethernets:
		  eth0:
		    match:
		       macaddress: "AA:FF:00:00:00:01"
		    addresses:
		      - 169.254.0.1/16
		  eth1:
		    match:
		      macaddress: "EE:00:00:00:00:__MAC_OCTET__"
		    addresses:
		      - __INSTANCE_IP__/24
		    gateway4: __GATEWAY__
		    nameservers:
		      addresses: [ 8.8.4.4, 8.8.8.8 ]
	*/

	/*
		# Boot configuration
		instance_ip=$VMS_NETWORK_PREFIX"."$(( $instance_number + 1 ))
		network_config_base64=$( \
			cat conf/cloud-init/network_config.yaml | \
			./tmpl.sh __INSTANCE_IP__ $instance_ip | \
			./tmpl.sh __MAC_OCTET__ $mac_octet | \
			./tmpl.sh __GATEWAY__ $VMS_NETWORK_PREFIX".1" | \
			gzip --stdout - | \
			base64 -w 0
		)
	*/
	nc := NetworkConfig{
		Version: 2,
		Ethernets: map[string]EthernetConfig{
			"eth0": {
				Match: EthernetMatcher{
					Macaddress: "AA:FF:00:00:00:01",
				},
				Addresses: []string{
					"169.254.0.1/16",
				},
				Gateway4:    "",
				Nameservers: nil,
			},
			"eth1": {
				Match: EthernetMatcher{
					Macaddress: "EE:00:00:00:00:01", // __MAC_OCTET__
				},
				Addresses: []string{
					"10.68.0.2/24", // __INSTANCE_IP__/24
				},
				Gateway4: "10.68.0.1", // __GATEWAY__
				Nameservers: &Nameservers{
					Addresses: []string{
						"1.1.1.1",
						"8.8.8.8",
					},
				},
			},
		},
	}
	ncBytes, err := yaml.Marshal(nc)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	//// Setting the Header fields is optional.
	//zw.Name = "a-new-hope.txt"
	//zw.Comment = "an epic space opera by George Lucas"
	//zw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)

	_, err = zw.Write(ncBytes)
	if err != nil {
		return "", err
	}

	if err := zw.Close(); err != nil {
		return "", err
	}

	cfg := base64.URLEncoding.EncodeToString(buf.Bytes())
	return cfg, nil
}

func buildData(instanceID string) error {
	/*
		#cloud-config
		users:
		  - name: fc
		    gecos: Firecracker user
		    shell: /bin/bash
		    groups: sudo
		    sudo: ALL=(ALL) NOPASSWD:ALL
		    ssh_authorized_keys:
		      - __SSH_PUB_KEY__
	*/
	keys, err := getSSHPubKeys("tamalsaha")
	if err != nil {
		return err
	}
	userData := UserData{
		Users: []User{
			{
				Name: "default",
			},
			{
				Name:              "runner",
				Gecos:             "GitHub Action Runner",
				Shell:             "/bin/bash",
				Groups:            strings.Join([]string{"sudo", "docker"}, ", "), // groups: users, admin
				Sudo:              "ALL=(ALL) NOPASSWD:ALL",
				SSHAuthorizedKeys: keys,
			},
		},
	}
	udBytes, err := yaml.Marshal(userData)
	if err != nil {
		return err
	}

	md := Metadata{
		InstanceID:    instanceID,
		LocalHostname: "gh-runner",
	}
	mdBytes, err := yaml.Marshal(md)
	if err != nil {
		return err
	}

	lc := LatestConfig{
		MetaData: string(mdBytes),
		UserData: "#cloud-config" + "\n" + string(udBytes),
	}

	data, err := json.Marshal(lc)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// https://github.com/tamalsaha.keys
func getSSHPubKeys(ghUsernames ...string) ([]string, error) {
	var keys []string
	var buf bytes.Buffer
	for _, username := range ghUsernames {
		resp, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", username))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		buf.Reset()
		if _, err = io.Copy(&buf, resp.Body); err != nil {
			return nil, err
		}
		userKeys := strings.Split(strings.TrimSpace(buf.String()), "\n")
		keys = append(keys, userKeys...)
	}

	return keys, nil
}
