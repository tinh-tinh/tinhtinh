package frameworks

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupGinApp() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Simple GET route
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// JSON response route
	r.GET("/json", func(c *gin.Context) {
		data := CreateTestData()
		c.JSON(http.StatusOK, data)
	})

	// JSON request/response route
	r.POST("/json", func(c *gin.Context) {
		var data TestData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	// Path parameter route
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// Query parameter route
	r.GET("/query", func(c *gin.Context) {
		name := c.DefaultQuery("name", "")
		c.JSON(http.StatusOK, gin.H{"name": name})
	})

	return r
}

func BenchmarkGin_SimpleGET(b *testing.B) {
	handler := setupGinApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkGin_JSONResponse(b *testing.B) {
	handler := setupGinApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/json", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkGin_JSONRequestResponse(b *testing.B) {
	handler := setupGinApp()
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

func BenchmarkGin_PathParam(b *testing.B) {
	handler := setupGinApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/user/123", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkGin_QueryParam(b *testing.B) {
	handler := setupGinApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(handler, "GET", "/query?name=test", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkGin_WithMiddleware(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add middlewares
	r.Use(func(c *gin.Context) {
		c.Next()
	})
	r.Use(func(c *gin.Context) {
		c.Next()
	})
	r.Use(func(c *gin.Context) {
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := PerformRequest(r, "GET", "/", nil)
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkGin_ParallelRequests(b *testing.B) {
	handler := setupGinApp()

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
