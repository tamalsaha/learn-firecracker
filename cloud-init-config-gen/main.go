package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	keys, err := getSSHPubKeys("tamalsaha")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(keys))
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
