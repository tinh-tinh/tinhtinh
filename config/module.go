package config

import (
	"log"

	"github.com/tinh-tinh/tinhtinh/core"
)

const ENV core.Provide = "ConfigEnv"

func ForRoot[E any](path ...string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var lastValue *E
		path = append([]string{".env"}, path...)
		for _, v := range path {
			env, err := New[E](v)
			if err != nil {
				continue
			}
			lastValue = env
		}

		configModule := module.New(core.NewModuleOptions{})

		if lastValue == nil {
			log.Println("env not found")
			return configModule
		}

		configModule.NewProvider(core.ProviderOptions{
			Name:  ENV,
			Value: lastValue,
		})
		configModule.Export(ENV)

		return configModule
	}
}
