package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

var (
	requestCounts = map[string]int{}
	mu            sync.Mutex
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/health")
	w.WriteHeader(http.StatusOK)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/ready")
	w.WriteHeader(http.StatusOK)
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	// recordRequest("/payload")
	n := rand.Intn(20)
	fibonacci := fibonacciSequence(n)
	response := map[string]interface{}{
		"number":    n,
		"fibonacci": fibonacci,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/metrics")
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestCounts)
}

func recordRequest(endpoint string) {
	mu.Lock()
	defer mu.Unlock()
	requestCounts[endpoint]++
}

func fibonacciSequence(n int) []int {

	fib := make([]int, n)

	for i := 0; i <= n-1; i++ {
		if i == 0 || i == 1 {
			fib[i] = i
		} else {
			fib[i] = fib[i-1] + fib[i-2]
		}
	}

	return fib
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/payload", payloadHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Println("Service is starting...")
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
