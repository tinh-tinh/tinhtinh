package main

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/microservices"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func appController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("test")

	ctrl.Get("/", func(ctx core.Ctx) error {
		client := microservices.Inject(module)
		if client == nil {
			return ctx.JSON(core.Map{"error": "client not found"})
		}
		defer client.Close()
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
	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{microservices.RegisterClient(microservices.ClientOptions{
				Addr: "localhost:8080",
			})},
			Controllers: []core.Controller{appController},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("api")

	app.Listen(3000)
}
