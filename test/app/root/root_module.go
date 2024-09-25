package root

import "github.com/tinh-tinh/tinhtinh/core"

func NewModule(module *core.DynamicModule) *core.DynamicModule {
	rootModule := module.New(core.NewModuleOptions{
		Scope:     core.Request,
		Providers: []core.Provider{NewProvider},
		Exports:   []core.Provider{NewProvider},
	})
	return rootModule
}
