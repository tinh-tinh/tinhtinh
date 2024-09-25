package user

import "github.com/tinh-tinh/tinhtinh/core"

func UserController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("user").Version("1")

	ctrl.Get("/test", func(ctx core.Ctx) {
		name := ctrl.Inject("user")
		ctx.JSON(core.Map{
			"data": name,
		})
	})

	ctrl.Get("/tip", func(ctx core.Ctx) {
		name := ctrl.Inject("user")
		ctx.JSON(core.Map{
			"data": name,
		})
	})

	ctrl.Get("/track", func(ctx core.Ctx) {
		name := ctrl.Inject("user")
		ctx.JSON(core.Map{
			"data": name,
		})
	})
	return ctrl
}
