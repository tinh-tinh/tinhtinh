package main

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/static"
)

func main() {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) {
			ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Imports:     []core.Module{static.ForRoot("upload")},
			Controllers: []core.Controller{appController},
		})
	}

	app := core.CreateFactory(appModule)
	app.Listen(8080)
}
