package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
)

func authController(module *core.DynamicModule) *core.DynamicController {
	authCtrl := module.NewController("Auth").Tag("Global")

	authCtrl.Pipe(
		core.Body(&dto.SignUpUser{}),
	).Post("/", func(ctx core.Ctx) {
		payload := ctx.Get(core.InBody).(*dto.SignUpUser)

		userService := authCtrl.Inject(USER_SERVICE).(CrudService)
		err := userService.Create(*payload)
		if err != nil {
			ctx.JSON(core.Map{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(core.Map{
			"status": "ok",
		})
	})

	return authCtrl
}
