package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, World!")
		})

		api.GET("/json", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello from Gin",
				"status":  "ok",
			})
		})

		api.POST("/json", func(c *gin.Context) {
			var data map[string]interface{}
			if err := c.BindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, data)
		})

		api.GET("/user/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "name": "User " + id})
		})
	}

	log.Println("Gin server starting on :3001")
	log.Fatal(r.Run(":3001"))
}
