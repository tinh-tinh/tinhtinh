package frameworks

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
)

func setupEchoApp() http.Handler {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Simple GET route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// JSON response route
	e.GET("/json", func(c echo.Context) error {
		data := CreateTestData()
		return c.JSON(http.StatusOK, data)
	})

	// JSON request/response route
	e.POST("/json", func(c echo.Context) error {
		var data TestData
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, data)
	})

	// Path parameter route
	e.GET("/user/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, map[string]string{"id": id})
	})

	// Query parameter route
	e.GET("/query", func(c echo.Context) error {
		name := c.QueryParam("name")
		return c.JSON(http.StatusOK, map[string]string{"name": name})
	})

	return e
}

func BenchmarkEcho_SimpleGET(b *testing.B) {
	handler := setupEchoApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkEcho_JSONResponse(b *testing.B) {
	handler := setupEchoApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/json", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkEcho_JSONRequestResponse(b *testing.B) {
	handler := setupEchoApp()
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

func BenchmarkEcho_PathParam(b *testing.B) {
	handler := setupEchoApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/user/123", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkEcho_QueryParam(b *testing.B) {
	handler := setupEchoApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/query?name=test", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkEcho_WithMiddleware(b *testing.B) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Add middlewares
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	})
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	})
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(e, "GET", "/", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkEcho_ParallelRequests(b *testing.B) {
	handler := setupEchoApp()

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
