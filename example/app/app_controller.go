package app

import "github.com/tinh-tinh/tinhtinh/api"

func NewController() *api.Controller {
	userController := api.NewController("users")

	userController.Get("/", func(ctx api.Ctx) {
		userService := NewService()
		data := userService.GetAll()
		ctx.JSON(api.Map{"data": data})
	})

	return userController
}
