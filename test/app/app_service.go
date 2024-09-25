package app

import "github.com/tinh-tinh/tinhtinh/core"

func AppService(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name:  "app",
		Value: "test",
	})
	return provider
}
