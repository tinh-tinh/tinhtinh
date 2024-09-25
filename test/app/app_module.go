package app

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/test/app/abc"
	"github.com/tinh-tinh/tinhtinh/test/app/root"
	"github.com/tinh-tinh/tinhtinh/test/app/user"
)

func NewModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			root.NewModule,
			abc.NewModule,
			user.NewModule,
		},
		Controllers: []core.Controller{AppController},
		Providers:   []core.Provider{AppService},
	})

	return appModule
}
