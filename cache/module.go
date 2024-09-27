package cache

import (
	"time"

	"github.com/tinh-tinh/tinhtinh/core"
)

type Store interface {
	Get(key string) interface{}
	Set(key string, val interface{}, ttl ...time.Duration)
	Keys() []string
	Count() int
	Has(key string) bool
	Delete(key string)
	Clear()
}

type Options struct {
	Store Store
	Ttl   time.Duration
	Max   int
}

const CACHE_MANAGER core.Provide = "CACHE_MANAGER"

func Register(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		memory := NewInMemory(MemoryOptions{
			Ttl: opt.Ttl,
			Max: opt.Max,
		})

		cacheModule := module.New(core.NewModuleOptions{})
		cacheModule.NewProvider(core.ProviderOptions{
			Name:  CACHE_MANAGER,
			Value: memory,
		})
		cacheModule.Export(CACHE_MANAGER)

		return cacheModule
	}
}
