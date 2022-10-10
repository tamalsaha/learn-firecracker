package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"sigs.k8s.io/yaml"
)

func main__() {
	keys, err := getSSHPubKeys("tamalsaha")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(keys))
}

func main() {
	err := buildData("i123")
	if err != nil {
		panic(err)
	}
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
