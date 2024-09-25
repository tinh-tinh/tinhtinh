package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

func NewModule(module *core.DynamicModule) *core.DynamicModule {
	userModule := module.New(core.NewModuleOptions{
		Scope:       core.Request,
		Controllers: []core.Controller{UserController, UserV2Controller},
		Providers:   []core.Provider{UserProvider},
	})

	return userModule
}
