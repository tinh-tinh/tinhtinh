package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/database/sql"
)

// func Module() *api.Module {
// 	userModule := api.NewModule(api.NewModuleOptions{
// 		Controllers: []api.ControllerParam{managerController, authController},
// 		Providers:   []api.ProviderParam{service},
// 	})

// 	return userModule
// }

func Module() *core.DynamicModule {
	userModule := core.NewModule(core.NewModuleOptions{
		Controllers: []core.Controller{managerController, authController},
		Providers: []core.Provider{
			sql.RegistryModel[User](),
			userService,
		},
	})

	return userModule
}
