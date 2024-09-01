package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
)

// func managerController(module *api.Module) *api.Controller {
// 	ctrl := api.NewController("users", module)

// 	ctrl.Pipe(
// 		api.Query[dto.FindUser](),
// 	).Get("/", func(ctx api.Ctx) {
// 		userService := ctrl.Inject("USER").(Service)
// 		data := userService.GetAll()
// 		ctx.JSON(api.Map{"data": data})
// 	})

// 	ctrl.Pipe(
// 		api.Body[dto.SignUpUser](),
// 	).Post("/", func(ctx api.Ctx) {
// 		ctx.JSON(api.Map{"data": "ok"})
// 	})

// 	return ctrl
// }

func managerController(module *core.DynamicModule) *core.DynamicController {
	ctrl := core.NewController("users", module)

	ctrl.Pipe(
		core.Query[dto.FindUser](),
	).Get("/", func(ctx core.Ctx) {
		userService := ctrl.Inject(USER_SERVICE).(Service)
		data := userService.GetAll()
		ctx.JSON(core.Map{"data": data})
	})

	ctrl.Pipe(
		core.Body[dto.SignUpUser](),
	).Post("/", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}
