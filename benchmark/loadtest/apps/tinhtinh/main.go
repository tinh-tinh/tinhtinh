package main

import (
	"log"
	"net/http"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func main() {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("api")

		ctrl.Get("", func(ctx core.Ctx) error {
			ctx.Res().Write([]byte("Hello, World!"))
			return nil
		})

		ctrl.Get("json", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"message": "Hello from Tinh Tinh",
				"status":  "ok",
			})
		})

		ctrl.Post("json", func(ctx core.Ctx) error {
			var data map[string]interface{}
			if err := ctx.BodyParser(&data); err != nil {
				return ctx.Status(http.StatusBadRequest).JSON(core.Map{"error": err.Error()})
			}
			return ctx.JSON(data)
		})

		ctrl.Get("user/:id", func(ctx core.Ctx) error {
			id := ctx.Path("id")
			return ctx.JSON(core.Map{"id": id, "name": "User " + id})
		})

		return ctrl
	}

	module := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	log.Println("Tinh Tinh server starting on :3000")
	app.Listen(3000)
}
