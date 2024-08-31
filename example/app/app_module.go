package app

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/example/app/user"
)

func NewModule() *api.Module {
	appModule := api.NewModule(api.NewModuleOptions{
		Import: []api.ModuleParam{user.Module},
	})

	return appModule
}
