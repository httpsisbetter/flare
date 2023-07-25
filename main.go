package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const outputFileName = "output.txt"

func main() {
	http.HandleFunc("/receive", handleIPs)
	fmt.Println("Go server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleIPs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var ipsRequest struct {
		IPs []string `json:"ips"`
	}

	err := json.NewDecoder(r.Body).Decode(&ipsRequest)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		http.Error(w, "Failed to create or open output file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	for _, ip := range ipsRequest.IPs {
		_, err := file.WriteString(ip + "\n")
		if err != nil {
			http.Error(w, "Failed to write IP to file", http.StatusInternalServerError)
			return
		}
	}

	response := map[string]string{"message": "IPs received and saved successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
