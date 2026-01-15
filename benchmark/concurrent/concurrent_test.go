package concurrent

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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

// BenchmarkConcurrent_10 tests with 10 concurrent requests
func BenchmarkConcurrent_10(b *testing.B) {
	handler := setupApp()
	concurrency := 10

	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	requestsPerGoroutine := b.N / concurrency

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					b.Errorf("Expected status 200, got %d", w.Code)
				}
			}
		}()
	}

	wg.Wait()
}

// BenchmarkConcurrent_100 tests with 100 concurrent requests
func BenchmarkConcurrent_100(b *testing.B) {
	handler := setupApp()
	concurrency := 100

	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	requestsPerGoroutine := b.N / concurrency

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					b.Errorf("Expected status 200, got %d", w.Code)
				}
			}
		}()
	}

	wg.Wait()
}

// BenchmarkConcurrent_1000 tests with 1000 concurrent requests
func BenchmarkConcurrent_1000(b *testing.B) {
	handler := setupApp()
	concurrency := 1000

	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	requestsPerGoroutine := b.N / concurrency

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					b.Errorf("Expected status 200, got %d", w.Code)
				}
			}
		}()
	}

	wg.Wait()
}

// BenchmarkConcurrent_10000 tests with 10000 concurrent requests
func BenchmarkConcurrent_10000(b *testing.B) {
	handler := setupApp()
	concurrency := 10000

	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	requestsPerGoroutine := b.N / concurrency
	if requestsPerGoroutine == 0 {
		requestsPerGoroutine = 1
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					b.Errorf("Expected status 200, got %d", w.Code)
				}
			}
		}()
	}

	wg.Wait()
}

// BenchmarkConcurrent_RunParallel uses Go's built-in parallel benchmark
func BenchmarkConcurrent_RunParallel(b *testing.B) {
	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

// BenchmarkConcurrent_WithContention tests with shared state contention
func BenchmarkConcurrent_WithContention(b *testing.B) {
	handler := setupApp()
	var counter int64

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&counter, 1)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})

	b.Logf("Total requests: %d", atomic.LoadInt64(&counter))
}

// BenchmarkConcurrent_ContextPooling tests context pooling under load
func BenchmarkConcurrent_ContextPooling(b *testing.B) {
	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/json", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

// BenchmarkConcurrent_Sustained tests sustained concurrent load
func BenchmarkConcurrent_Sustained(b *testing.B) {
	handler := setupApp()
	duration := 5 * time.Second
	concurrency := 100

	b.ResetTimer()

	var wg sync.WaitGroup
	var totalRequests int64
	stop := make(chan struct{})

	// Start timer
	go func() {
		time.Sleep(duration)
		close(stop)
	}()

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					req := httptest.NewRequest("GET", "/", nil)
					w := httptest.NewRecorder()
					handler.ServeHTTP(w, req)
					atomic.AddInt64(&totalRequests, 1)
				}
			}
		}()
	}

	wg.Wait()

	total := atomic.LoadInt64(&totalRequests)
	rps := float64(total) / duration.Seconds()

	b.Logf("Duration: %v", duration)
	b.Logf("Concurrency: %d", concurrency)
	b.Logf("Total requests: %d", total)
	b.Logf("Requests per second: %.2f", rps)
}
