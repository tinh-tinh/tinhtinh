package token

import "github.com/tinh-tinh/tinhtinh/core"

const TOKEN core.Provide = "TOKEN"

func Register(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		provider := core.NewProvider(module)
		provider.Set(TOKEN, NewProvider(opt))

		return module
	}
}
