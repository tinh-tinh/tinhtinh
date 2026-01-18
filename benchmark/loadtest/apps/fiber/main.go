package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	api := app.Group("/api")
	{
		api.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})

		api.Get("/json", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Hello from Fiber",
				"status":  "ok",
			})
		})

		api.Post("/json", func(c *fiber.Ctx) error {
			var data map[string]interface{}
			if err := c.BodyParser(&data); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(data)
		})

		api.Get("/user/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			return c.JSON(fiber.Map{"id": id, "name": "User " + id})
		})
	}

	log.Println("Fiber server starting on :3003")
	log.Fatal(app.Listen(":3003"))
}
