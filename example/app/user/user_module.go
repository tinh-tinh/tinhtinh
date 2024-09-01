package user

import (
	"fmt"
	"os/user"

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

func Module(module *core.DynamicModule) *core.DynamicModule {
	fmt.Println(module.Ref(sql.ConnectDB))
	userModule := core.NewModule(core.NewModuleOptions{
		Imports:     []core.Module{sql.ForFeature(&user.User{})},
		Controllers: []core.Controller{managerController, authController},
		Providers:   []core.Provider{userService},
	})

	return userModule
}
