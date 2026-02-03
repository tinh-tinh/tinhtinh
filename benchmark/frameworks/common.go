package frameworks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestData represents common test data structure
type TestData struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// BenchmarkConfig holds configuration for benchmarks
type BenchmarkConfig struct {
	Routes      int
	Middlewares int
	PayloadSize int
}

// DefaultConfig returns default benchmark configuration
func DefaultConfig() BenchmarkConfig {
	return BenchmarkConfig{
		Routes:      10,
		Middlewares: 3,
		PayloadSize: 1024,
	}
}

// CreateTestData creates sample test data
func CreateTestData() TestData {
	return TestData{
		ID:      1,
		Name:    "John Doe",
		Email:   "john@example.com",
		Message: "This is a test message for benchmarking purposes",
	}
}

// CreateJSONPayload creates a JSON payload for testing
func CreateJSONPayload(data interface{}) []byte {
	payload, _ := json.Marshal(data)
	return payload
}

// PerformRequest performs an HTTP request for benchmarking
func PerformRequest(handler http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// BenchmarkResult holds benchmark results
type BenchmarkResult struct {
	Name        string
	NsPerOp     int64
	BytesPerOp  int64
	AllocsPerOp int64
	MBPerSec    float64
}

// CompareResults compares two benchmark results
func CompareResults(baseline, current BenchmarkResult) float64 {
	if baseline.NsPerOp == 0 {
		return 0
	}
	return float64(current.NsPerOp-baseline.NsPerOp) / float64(baseline.NsPerOp) * 100
}

// SimpleMiddleware is a simple middleware for testing
func SimpleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple middleware logic
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware simulates a logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate logging
		_ = r.Method + " " + r.URL.Path
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware simulates an auth middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate auth check
		token := r.Header.Get("Authorization")
		if token == "" {
			// For benchmarking, we don't actually fail
		}
		next.ServeHTTP(w, r)
	})
}

// Helper function to create request body
func CreateRequestBody(data interface{}) *bytes.Buffer {
	payload := CreateJSONPayload(data)
	return bytes.NewBuffer(payload)
}

// ValidateResponse validates the response
func ValidateResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	if w.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, w.Code)
	}
}

// ParseJSONResponse parses JSON response
func ParseJSONResponse(w *httptest.ResponseRecorder, v interface{}) error {
	return json.Unmarshal(w.Body.Bytes(), v)
}
