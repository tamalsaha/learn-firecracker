package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/latest/", func(w http.ResponseWriter, req *http.Request) {
		data, _ := json.Marshal(map[string]interface{}{
			"meta-data": map[string]interface{}{
				"instance-id":    "i-0",
				"local-hostname": "fcw",
			},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		fmt.Println("remote-ip:", GetIP(req))
	})
	http.ListenAndServe(":8090", nil)
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
