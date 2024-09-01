package post

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/database/sql"
)

func Module(module *core.DynamicModule) *core.DynamicModule {
	postModule := module.New(core.NewModuleOptions{
		Imports:     []core.Module{sql.ForFeature(&Post{})},
		Controllers: []core.Controller{controller},
		Providers:   []core.Provider{service},
	})

	return postModule
}
