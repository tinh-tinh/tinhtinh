package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	api := e.Group("/api")
	{
		api.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		api.GET("/json", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Hello from Echo",
				"status":  "ok",
			})
		})

		api.POST("/json", func(c echo.Context) error {
			var data map[string]interface{}
			if err := c.Bind(&data); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			}
			return c.JSON(http.StatusOK, data)
		})

		api.GET("/user/:id", func(c echo.Context) error {
			id := c.Param("id")
			return c.JSON(http.StatusOK, map[string]interface{}{"id": id, "name": "User " + id})
		})
	}

	log.Println("Echo server starting on :3002")
	log.Fatal(e.Start(":3002"))
}
