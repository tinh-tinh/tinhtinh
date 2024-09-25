package main

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/cors"
	"github.com/tinh-tinh/tinhtinh/middleware/helmet"
	"github.com/tinh-tinh/tinhtinh/middleware/logger"
	"github.com/tinh-tinh/tinhtinh/test/app"
)

func main() {
	app := core.CreateFactory(app.NewModule, "api").EnableCors(cors.CorsOptions{
		AllowedMethods: []string{"POST", "GET"},
		AllowedHeaders: []string{"*"},
	})

	app.EnableVersioning(core.VersionOptions{
		Type: core.MediaTypeVersion,
		Key:  "v=",
	})

	app.Use(logger.Middleware(logger.MiddlewareOptions{
		Rotate: true,
		Format: "${method} ${path} ${status} ${latency}",
	}))

	h := helmet.New(helmet.HelmetOptions{
		XPoweredBy: helmet.XPoweredBy{
			Enabled: true, Value: "tinhtinh"},
	})
	app.Use(h.Handler)
	app.BeforeShutdown(func() {
		fmt.Println("Before shutdown")
	})
	app.AfterShutdown(func() {
		fmt.Print("After shutdown")
	})
	app.Listen(3000)
}
