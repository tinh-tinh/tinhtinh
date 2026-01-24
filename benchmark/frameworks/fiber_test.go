package frameworks

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func BenchmarkFiber_SimpleGET(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func BenchmarkFiber_JSONResponse(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/json", func(c *fiber.Ctx) error {
		data := CreateTestData()
		return c.JSON(data)
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/json", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func BenchmarkFiber_JSONRequestResponse(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/json", func(c *fiber.Ctx) error {
		var data TestData
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	data := CreateTestData()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		body := CreateRequestBody(data)
		req, _ := http.NewRequest("POST", "/json", body)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func BenchmarkFiber_PathParam(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{"id": id})
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/user/123", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func BenchmarkFiber_QueryParam(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/query", func(c *fiber.Ctx) error {
		name := c.Query("name", "")
		return c.JSON(fiber.Map{"name": name})
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/query?name=test", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func BenchmarkFiber_WithMiddleware(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// Add middlewares
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func BenchmarkFiber_ParallelRequests(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/", nil)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("Expected status 200, got %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}
