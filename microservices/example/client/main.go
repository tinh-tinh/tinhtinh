package main

import (
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func appController(module core.Module) core.Controller {
	ctrl := module.NewController("test")

	ctrl.Get("update", func(ctx core.Ctx) error {
		client := microservices.Inject(module)
		if client == nil {
			return ctx.JSON(core.Map{"error": "client not found"})
		}
		// Example JSON messages to send
		messages := []User{
			{"haha", 30},
			{"hihi", 25},
			{"huhu", 35},
		}

		for _, msg := range messages {
			client.Send("user.updated", msg)
		}

		return ctx.JSON(core.Map{"data": "update"})
	})

	ctrl.Get("", func(ctx core.Ctx) error {
		client := microservices.Inject(module)
		if client == nil {
			return ctx.JSON(core.Map{"error": "client not found"})
		}
		// Example JSON messages to send
		messages := []User{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 35},
		}

		for _, msg := range messages {
			client.Send("user.created", msg)
		}

		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}

func main() {
	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:     []core.Modules{microservices.RegisterClient("localhost:8080")},
			Controllers: []core.Controllers{appController},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("api")

	app.Listen(3000)
}
