package main

import "github.com/tinh-tinh/tinhtinh/core"

func main() {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "1",
			})
		})

		ctrl.Post("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})

		ctrl.Patch("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "3",
			})
		})

		ctrl.Put("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "4",
			})
		})

		ctrl.Delete("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "5",
			})
		})
		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	app.Listen(5000)
}
