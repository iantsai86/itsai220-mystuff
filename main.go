package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microservice_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"endpoint"},
	)

	mu sync.Mutex
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(requestCounts)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/health")
	w.WriteHeader(http.StatusOK)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/ready")
	w.WriteHeader(http.StatusOK)
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	recordRequest("/payload")
	n := rand.Intn(20) + 1
	fibonacci := fibonacciSequence(n)
	response := map[string]interface{}{
		"number":    n,
		"fibonacci": fibonacci,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func recordRequest(endpoint string) {
	requestCounts.WithLabelValues(endpoint).Inc()
}

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
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP) // Prometheus metrics endpoint

	log.Println("Service is starting...")
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
