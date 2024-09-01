package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

// func Module() *api.Module {
// 	userModule := api.NewModule(api.NewModuleOptions{
// 		Controllers: []api.ControllerParam{managerController, authController},
// 		Providers:   []api.ProviderParam{service},
// 	})

// 	return userModule
// }

func Module(m *core.DynamicModule) *core.DynamicModule {
	userModule := core.NewModule(core.NewModuleOptions{
		Controllers: []core.Controller{managerController, authController},
		Providers:   []core.Provider{userService},
	})

	return userModule
}
