package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/latest/", func(w http.ResponseWriter, req *http.Request) {
		data, _ := json.Marshal(map[string]interface{}{
			"meta-data": map[string]interface{}{
				"instance-id": "i-0",
				"hostname":    "fcw",
			},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	http.ListenAndServe(":8090", nil)
}
