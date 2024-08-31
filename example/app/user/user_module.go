package user

import (
	"github.com/tinh-tinh/tinhtinh/api"
)

func Module() *api.Module {
	userModule := api.NewModule(api.NewModuleOptions{
		Controllers: []api.ControllerParam{managerController, authController},
		Providers:   []api.ProviderParam{service},
	})

	return userModule
}
