package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
	req, err := http.NewRequest("GET", "/payload", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(payloadHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	number, ok := response["number"].(float64)
	if !ok {
		t.Errorf("Response body does not contain a valid number field")
	}

	fibonacci, ok := response["fibonacci"].([]interface{})
	if !ok {
		t.Errorf("Response body does not contain a valid fibonacci field")
	}

	// Check if the Fibonacci sequence length is equal to the number
	if len(fibonacci) != int(number) {
		t.Errorf("Fibonacci sequence length %v does not match the number %v", len(fibonacci), number)
	}
}

// func TestMetricsHandler(t *testing.T) {
// 	// Make a request to the /payload endpoint to ensure there are some metrics
// 	req, err := http.NewRequest("GET", "/payload", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(payloadHandler)

// 	handler.ServeHTTP(rr, req)

// 	// Now make a request to the /metrics endpoint
// 	req, err = http.NewRequest("GET", "/metrics", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr = httptest.NewRecorder()
// 	handler = http.HandlerFunc(metricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	var metrics map[string]interface{}
// 	err = json.Unmarshal(rr.Body.Bytes(), &metrics)
// 	if err != nil {
// 		t.Fatalf("Failed to unmarshal metrics response body: %v", err)
// 	}

// 	// Verify that the metrics contain the expected endpoints
// 	if _, ok := metrics["/health"]; !ok {
// 		t.Error("Metrics response does not contain /health request count")
// 	}
// 	if _, ok := metrics["/ready"]; !ok {
// 		t.Error("Metrics response does not contain /ready request count")
// 	}
// 	if _, ok := metrics["/payload"]; !ok {
// 		t.Error("Metrics response does not contain /payload request count")
// 	}
// 	if _, ok := metrics["/metrics"]; !ok {
// 		t.Error("Metrics response does not contain /metrics request count")
// 	}
// }
