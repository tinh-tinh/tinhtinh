package app

import "github.com/tinh-tinh/tinhtinh/core"

func AppController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("test")

	ctrl.Get("/", func(ctx core.Ctx) {
		name := ctrl.Inject("sjkfkjgjkebjgrkebjkgb")
		ctx.JSON(core.Map{
			"data": name,
		})
	})

	return ctrl
}
