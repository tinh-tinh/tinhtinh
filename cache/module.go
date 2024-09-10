package cache

import (
	"time"

	"github.com/tinh-tinh/tinhtinh/core"
)

type Store[K any, V any] interface {
	Get(key K) *V
	Set(key K, val V, ttl ...time.Duration)
	Keys() []K
	Count() int
	Has(key K) bool
	Delete(key K)
	Clear()
}

type Options[K string, V interface{}] struct {
	Store Store[K, V]
	Ttl   time.Duration
	Max   int
}

const CACHE_MANAGER core.Provide = "CACHE_MANAGER"

func Register[K string, V any](opt Options[K, V]) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		memory := NewInMemory[K, V](MemoryOptions{
			Ttl: opt.Ttl,
			Max: opt.Max,
		})

		cacheModule := module.New(core.NewModuleOptions{})
		cacheModule.NewProvider(memory, CACHE_MANAGER)
		cacheModule.Export(CACHE_MANAGER)

		return cacheModule
	}
}
