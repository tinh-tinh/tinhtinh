package user

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
)

func NewController() *api.Controller {
	userController := api.NewController("users")

	userController.Pipe(
		api.Query[dto.FindUser](),
	).Get("/", func(ctx api.Ctx) {
		userService := NewService()
		data := userService.GetAll()
		ctx.JSON(api.Map{"data": data})
	})

	userController.Pipe(
		api.Body[dto.SignUpUser](),
	).Post("/", func(ctx api.Ctx) {
		ctx.JSON(api.Map{"data": "ok"})
	})

	return userController
}
