package profiling

import (
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"

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

		ctrl.Post("json", func(ctx core.Ctx) error {
			var data map[string]interface{}
			if err := ctx.BodyParser(&data); err != nil {
				return ctx.JSON(core.Map{"error": err.Error()})
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

// BenchmarkCPUProfile benchmarks CPU usage
func BenchmarkCPUProfile(b *testing.B) {
	// Create CPU profile
	f, err := os.Create("cpu.prof")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		b.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// BenchmarkMemoryProfile benchmarks memory allocations
func BenchmarkMemoryProfile(b *testing.B) {
	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}

	b.StopTimer()

	// Write memory profile
	f, err := os.Create("mem.prof")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	runtime.GC() // Force GC before writing profile
	if err := pprof.WriteHeapProfile(f); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkGoroutineProfile benchmarks goroutine creation
func BenchmarkGoroutineProfile(b *testing.B) {
	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}

	b.StopTimer()

	// Write goroutine profile
	f, err := os.Create("goroutine.prof")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkAllocations benchmarks memory allocations
func BenchmarkAllocations(b *testing.B) {
	handler := setupApp()

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	b.StopTimer()
	runtime.ReadMemStats(&m2)

	// Report allocation statistics
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "B/op")
	b.ReportMetric(float64(m2.Mallocs-m1.Mallocs)/float64(b.N), "allocs/op")
}

// BenchmarkContextPooling benchmarks context pooling efficiency
func BenchmarkContextPooling(b *testing.B) {
	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		}
	})
}

// BenchmarkJSONSerialization benchmarks JSON encoding/decoding
func BenchmarkJSONSerialization(b *testing.B) {
	handler := setupApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}
