package tcp_test

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthApp(addr string) *core.App {
	const USER_SERVICE = "user_service"
	type UserService struct {
		users []*User
		total int
	}

	service := func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: USER_SERVICE,
			Value: &UserService{
				users: make([]*User, 0),
				total: 0,
			},
		})

		return prd
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("auth")

		svc := module.Ref(USER_SERVICE).(*UserService)
		ctrl.Post("/login", func(ctx core.Ctx) error {
			var user *User
			input := ctx.BodyParser(&user)

			client := microservices.Inject(module)

		})

		return ctrl
	}
}
