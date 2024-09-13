package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestReadyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestPayloadHandler(t *testing.T) {
	// Reset the metrics before testing
	requestCounts.Reset()

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/payload", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test HTTP handler
	handler := http.HandlerFunc(payloadHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Validate the response body
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if _, ok := response["number"]; !ok {
		t.Errorf("Response body missing 'number' field")
	}
	if _, ok := response["fibonacci"]; !ok {
		t.Errorf("Response body missing 'fibonacci' field")
	}

	// Validate the counter metric
	expected := `
# HELP service_requests_total Total number of requests received
# TYPE service_requests_total counter
service_requests_total{endpoint="/payload"} 1
`
	if err := testutil.GatherAndCompare(prometheus.DefaultGatherer, strings.NewReader(expected), "service_requests_total"); err != nil {
		t.Errorf("Unexpected metrics: %v", err)
	}
}
