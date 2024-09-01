package user

import (
	"os/user"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/database/sql"
)

func Module(module *core.DynamicModule) *core.DynamicModule {
	userModule := module.New(core.NewModuleOptions{
		Imports:     []core.Module{sql.ForFeature(&user.User{})},
		Controllers: []core.Controller{managerController, authController},
		Providers:   []core.Provider{service},
	})

	return userModule
}
