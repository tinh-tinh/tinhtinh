package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
	"github.com/tinh-tinh/tinhtinh/jwt"
)

func managerController(module *core.DynamicModule) *core.DynamicController {
	ctrl := core.NewController("Users", module)

	ctrl.GuardWithCtrl(jwt.Guard).Pipe(
		core.Query(&dto.FindUser{}),
	).Get("/", func(ctx core.Ctx) {
		userService := ctrl.Inject(USER_SERVICE).(CrudService)
		data := userService.GetAll()
		ctx.JSON(core.Map{"data": data})
	})

	ctrl.Pipe(
		core.Body(&dto.SignUpUser{}),
	).Post("/", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}
