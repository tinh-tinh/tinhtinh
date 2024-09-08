package token

import "github.com/tinh-tinh/tinhtinh/core"

const TOKEN core.Provide = "TOKEN"

func Register(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		tokenModule := module.New(core.NewModuleOptions{})
		tokenModule.NewProvider(NewJwt(opt), TOKEN)
		tokenModule.Export(TOKEN)

		return tokenModule
	}
}
