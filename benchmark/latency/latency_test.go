package latency

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

// setupApp creates a test application
func setupApp() http.Handler {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			ctx.Res().Write([]byte("Hello, World!"))
			return nil
		})

		ctrl.Get("json", func(ctx core.Ctx) error {
			data := map[string]interface{}{
				"id":      1,
				"name":    "Test User",
				"email":   "test@example.com",
				"message": "This is a test message",
			}
			return ctx.JSON(data)
		})

		return ctrl
	}

	module := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})
	}

	app := core.CreateFactory(module)
	return app.PrepareBeforeListen()
}

// BenchmarkLatency_SimpleGET measures latency for simple GET requests
func BenchmarkLatency_SimpleGET(b *testing.B) {
	handler := setupApp()
	durations := make([]float64, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		duration := time.Since(start)
		durations[i] = float64(duration.Nanoseconds())

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}

	b.StopTimer()

	// Calculate and report latency statistics
	stats := Calculate(durations)
	b.Logf("\n%s", GenerateReport(stats, "Simple GET Request Latency"))
	b.Logf("%s", GenerateHistogram(durations, 10))
}

// BenchmarkLatency_JSONResponse measures latency for JSON responses
func BenchmarkLatency_JSONResponse(b *testing.B) {
	handler := setupApp()
	durations := make([]float64, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		duration := time.Since(start)
		durations[i] = float64(duration.Nanoseconds())

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}

	b.StopTimer()

	// Calculate and report latency statistics
	stats := Calculate(durations)
	b.Logf("\n%s", GenerateReport(stats, "JSON Response Latency"))
	b.Logf("%s", GenerateHistogram(durations, 10))
}

// BenchmarkLatency_Concurrent measures latency under concurrent load
func BenchmarkLatency_Concurrent(b *testing.B) {
	handler := setupApp()

	b.RunParallel(func(pb *testing.PB) {
		var localDurations []float64

		for pb.Next() {
			start := time.Now()

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			duration := time.Since(start)
			localDurations = append(localDurations, float64(duration.Nanoseconds()))

			if w.Code != http.StatusOK {
				b.Fatalf("Expected status 200, got %d", w.Code)
			}
		}

		if len(localDurations) > 0 {
			stats := Calculate(localDurations)
			b.Logf("\n%s", GenerateReport(stats, "Concurrent Request Latency (per goroutine)"))
		}
	})
}

// BenchmarkLatency_PercentileAccuracy tests percentile calculation accuracy
func BenchmarkLatency_PercentileAccuracy(b *testing.B) {
	// Create test data with known distribution
	testData := []float64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	}

	stats := Calculate(testData)

	b.Logf("Test Data Statistics:")
	b.Logf("Min: %.2f", stats.Min)
	b.Logf("p50: %.2f (expected: ~10.5)", stats.P50)
	b.Logf("p95: %.2f (expected: ~19.05)", stats.P95)
	b.Logf("p99: %.2f (expected: ~19.81)", stats.P99)
	b.Logf("Max: %.2f", stats.Max)
	b.Logf("Mean: %.2f (expected: 10.5)", stats.Mean)
}
