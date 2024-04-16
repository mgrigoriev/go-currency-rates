package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func listenAndServe() {
	fmt.Println("Starting HTTP server at http://" + httpAddr)

	http.HandleFunc("/", rates)
	http.ListenAndServe(httpAddr, nil)
}

func rates(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(ratesCache)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
