package frameworks

import (
	"net/http"
	"testing"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

// setupTinhTinhApp creates a Tinh Tinh app for benchmarking
func setupTinhTinhApp() http.Handler {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			ctx.Res().Write([]byte("Hello, World!"))
			return nil
		})

		ctrl.Get("json", func(ctx core.Ctx) error {
			data := CreateTestData()
			return ctx.JSON(data)
		})

		ctrl.Post("json", func(ctx core.Ctx) error {
			var data TestData
			if err := ctx.BodyParser(&data); err != nil {
				return ctx.JSON(core.Map{"error": err.Error()})
			}
			return ctx.JSON(data)
		})

		ctrl.Get("user/:id", func(ctx core.Ctx) error {
			id := ctx.Path("id")
			return ctx.JSON(core.Map{"id": id})
		})

		ctrl.Get("query", func(ctx core.Ctx) error {
			name := ctx.Query("name")
			return ctx.JSON(core.Map{"name": name})
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

func BenchmarkTinhTinh_SimpleGET(b *testing.B) {
	handler := setupTinhTinhApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkTinhTinh_JSONResponse(b *testing.B) {
	handler := setupTinhTinhApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/json", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkTinhTinh_JSONRequestResponse(b *testing.B) {
	handler := setupTinhTinhApp()
	data := CreateTestData()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		body := CreateRequestBody(data)
		w := PerformRequest(handler, "POST", "/json", body)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkTinhTinh_PathParam(b *testing.B) {
	handler := setupTinhTinhApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/user/123", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkTinhTinh_QueryParam(b *testing.B) {
	handler := setupTinhTinhApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/query?name=test", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkTinhTinh_WithMiddleware(b *testing.B) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")
		ctrl.Get("", func(ctx core.Ctx) error {
			ctx.Res().Write([]byte("Hello, World!"))
			return nil
		})
		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		// Add middlewares
		appModule.Use(func(ctx core.Ctx) error {
			return ctx.Next()
		})
		appModule.Use(func(ctx core.Ctx) error {
			return ctx.Next()
		})
		appModule.Use(func(ctx core.Ctx) error {
			return ctx.Next()
		})

		return appModule
	}

	app := core.CreateFactory(module)
	handler := app.PrepareBeforeListen()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkTinhTinh_ParallelRequests(b *testing.B) {
	handler := setupTinhTinhApp()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := PerformRequest(handler, "GET", "/", nil)
			if w.Code != http.StatusOK {
				b.Fatalf("Expected status 200, got %d", w.Code)
			}
		}
	})
}
