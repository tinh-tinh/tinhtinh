package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/core"
)

type Module struct {
	sync.Pool
}

func Register[E any](path string) (*E, error) {
	if path == "" {
		path = ".env"
	}
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	var env E
	Scan(&env)
	return &env, nil
}

func GetRaw(key string) string {
	return os.Getenv(key)
}

const ENV core.Provide = "ConfigEnv"

func ForRoot[E any](path ...string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var lastValue *E
		path = append([]string{".env"}, path...)
		for _, v := range path {
			env, err := Register[E](v)
			if err != nil {
				continue
			}
			lastValue = env
		}

		if lastValue == nil {
			panic("env not found")
		}

		configModule := module.New(core.NewModuleOptions{})

		provider := core.NewProvider(configModule)
		provider.Set(ENV, *lastValue)
		provider.Export(ENV)

		return configModule
	}
}
