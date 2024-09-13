package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// requestCounts is a Prometheus CounterVec metric to track the number of requests
	// to different endpoints. It uses labels to distinguish between different endpoints.
	requestCounts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"endpoint"},
	)
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(requestCounts)
}

// healthHandler handles requests to the /health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/health")
	w.WriteHeader(http.StatusOK)
}

// readyHandler handles requests to the /ready endpoint
func readyHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/ready")
	w.WriteHeader(http.StatusOK)
}

// payloadHandler handles requests to the /payload endpoint
func payloadHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/payload")

	n := rand.Intn(20)
	fibonacci := fibonacciSequence(n)
	response := map[string]interface{}{
		"number":    n,
		"fibonacci": fibonacci,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// recordRequest increments the request counter for a given endpoint
func recordRequest(endpoint string) {
	requestCounts.WithLabelValues(endpoint).Inc()
}

// fibonacciSequence computes the Fibonacci sequence up to the nth number
func fibonacciSequence(n int) []int {
	if n <= 0 {
		return []int{}
	}
	seq := []int{0, 1}
	for i := 2; i < n; i++ {
		seq = append(seq, seq[i-1]+seq[i-2])
	}
	return seq
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/payload", payloadHandler)

	// Set up the /metrics endpoint for Prometheus to scrape metrics
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	log.Println("Service is starting...")
	// Create and start the HTTP server
	server := &http.Server{
		Addr:    ":8081",
		Handler: nil,
	}

	// Start the server and log any errors
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
