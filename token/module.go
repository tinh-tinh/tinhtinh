package token

import "github.com/tinh-tinh/tinhtinh/core"

const TOKEN core.Provide = "TOKEN"

func Register(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		tokenModule := module.New(core.NewModuleOptions{})

		provider := core.NewProvider(tokenModule)
		provider.Set(TOKEN, NewProvider(opt))
		provider.Export(TOKEN)

		return tokenModule
	}
}
