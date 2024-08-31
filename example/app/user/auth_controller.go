package user

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
)

func authController(module *api.Module) *api.Controller {
	authCtrl := api.NewController("auth", module)

	authCtrl.Pipe(api.Body[dto.SignUpUser]()).Post("/", func(ctx api.Ctx) {
		payload := ctx.Get(api.Payload).(dto.SignUpUser)

		userService := authCtrl.Inject("USER").(Service)
		err := userService.Create(payload)
		if err != nil {
			ctx.JSON(api.Map{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(api.Map{
			"status": "ok",
		})
	})

	return authCtrl
}
