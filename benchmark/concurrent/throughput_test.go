package concurrent

import (
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkThroughput_SimpleGET measures requests per second for simple GET
func BenchmarkThroughput_SimpleGET(b *testing.B) {
	handler := setupApp()

	start := time.Now()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	duration := time.Since(start)
	rps := float64(b.N) / duration.Seconds()

	b.ReportMetric(rps, "req/s")
	b.Logf("Throughput: %.2f requests/second", rps)
}

// BenchmarkThroughput_JSONResponse measures RPS for JSON responses
func BenchmarkThroughput_JSONResponse(b *testing.B) {
	handler := setupApp()

	start := time.Now()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	duration := time.Since(start)
	rps := float64(b.N) / duration.Seconds()

	b.ReportMetric(rps, "req/s")
	b.Logf("Throughput: %.2f requests/second", rps)
}

// BenchmarkThroughput_Parallel measures parallel throughput
func BenchmarkThroughput_Parallel(b *testing.B) {
	handler := setupApp()
	var counter int64

	start := time.Now()
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			atomic.AddInt64(&counter, 1)
		}
	})

	duration := time.Since(start)
	total := atomic.LoadInt64(&counter)
	rps := float64(total) / duration.Seconds()

	b.ReportMetric(rps, "req/s")
	b.Logf("Throughput: %.2f requests/second", rps)
	b.Logf("Total requests: %d", total)
}

// BenchmarkThroughput_KeepAlive measures throughput with keep-alive
func BenchmarkThroughput_KeepAlive(b *testing.B) {
	handler := setupApp()

	start := time.Now()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Connection", "keep-alive")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	duration := time.Since(start)
	rps := float64(b.N) / duration.Seconds()

	b.ReportMetric(rps, "req/s")
	b.Logf("Throughput (keep-alive): %.2f requests/second", rps)
}

// BenchmarkThroughput_Sustained measures sustained throughput over time
func BenchmarkThroughput_Sustained(b *testing.B) {
	handler := setupApp()
	duration := 10 * time.Second

	var totalRequests int64
	stop := make(chan struct{})

	// Start timer
	go func() {
		time.Sleep(duration)
		close(stop)
	}()

	b.ResetTimer()

	start := time.Now()

	// Single goroutine sustained load
	for {
		select {
		case <-stop:
			goto done
		default:
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			atomic.AddInt64(&totalRequests, 1)
		}
	}

done:
	elapsed := time.Since(start)
	total := atomic.LoadInt64(&totalRequests)
	rps := float64(total) / elapsed.Seconds()

	b.Logf("Duration: %v", elapsed)
	b.Logf("Total requests: %d", total)
	b.Logf("Sustained throughput: %.2f requests/second", rps)
	b.ReportMetric(rps, "req/s")
}

// BenchmarkThroughput_BurstLoad measures throughput under burst load
func BenchmarkThroughput_BurstLoad(b *testing.B) {
	handler := setupApp()
	burstSize := 1000

	b.ResetTimer()
	b.ReportAllocs()

	start := time.Now()

	for i := 0; i < b.N/burstSize; i++ {
		// Send burst of requests
		for j := 0; j < burstSize; j++ {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		}

		// Small pause between bursts
		time.Sleep(10 * time.Millisecond)
	}

	duration := time.Since(start)
	rps := float64(b.N) / duration.Seconds()

	b.ReportMetric(rps, "req/s")
	b.Logf("Burst throughput: %.2f requests/second", rps)
	b.Logf("Burst size: %d", burstSize)
}

// BenchmarkThroughput_ConnectionPooling measures connection pooling efficiency
func BenchmarkThroughput_ConnectionPooling(b *testing.B) {
	handler := setupApp()

	// Simulate connection pooling by reusing request objects
	req := httptest.NewRequest("GET", "/", nil)

	start := time.Now()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	duration := time.Since(start)
	rps := float64(b.N) / duration.Seconds()

	b.ReportMetric(rps, "req/s")
	b.Logf("Throughput (pooled): %.2f requests/second", rps)
}
