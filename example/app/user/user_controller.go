package user

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
)

func managerController(module *api.Module) *api.Controller {
	ctrl := api.NewController("users", module)

	ctrl.Pipe(
		api.Query[dto.FindUser](),
	).Get("/", func(ctx api.Ctx) {
		userService := ctrl.Inject("USER").(Service)
		data := userService.GetAll()
		ctx.JSON(api.Map{"data": data})
	})

	ctrl.Pipe(
		api.Body[dto.SignUpUser](),
	).Post("/", func(ctx api.Ctx) {
		ctx.JSON(api.Map{"data": "ok"})
	})

	return ctrl
}
