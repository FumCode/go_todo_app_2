package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "ok",
		Message: "Api is working",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func main() {

	http.HandleFunc("/health", healthHandler)

	fmt.Println("App is running in port 3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting app", err)
	}
}
