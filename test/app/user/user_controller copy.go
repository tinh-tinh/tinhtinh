package user

import "github.com/tinh-tinh/tinhtinh/core"

func UserV2Controller(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("user").Version("2")

	ctrl.Get("/test", func(ctx core.Ctx) {
		ctx.JSON(core.Map{
			"data": "2",
		})
	})

	ctrl.Get("/tip", func(ctx core.Ctx) {
		ctx.JSON(core.Map{
			"data": "2",
		})
	})

	ctrl.Get("/track", func(ctx core.Ctx) {
		ctx.JSON(core.Map{
			"data": "2",
		})
	})
	return ctrl
}
