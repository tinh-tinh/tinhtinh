package app

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app/user"
)

// func NewModule() *api.Module {
// 	appModule := api.NewModule(api.NewModuleOptions{
// 		Import: []api.ModuleParam{user.Module},
// 	})

// 	return appModule
// }

func NewModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{user.Module},
	})

	return appModule
}
