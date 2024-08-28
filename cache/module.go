package cache

import (
	"sync"
	"time"
)

// type CacheType string

// const (
// 	InMemory CacheType = "memory"
// 	Redis    CacheType = "redis"
// )

// type ModuleOptions struct {
// 	Type CacheType
// }

// type Module struct {
// 	memory *Memory
// }

var pool sync.Pool

func Register() {
	pool = sync.Pool{
		New: func() any {
			return NewInMemory()
		},
	}
}

func Get(key string) interface{} {
	m := pool.Get().(*Memory)
	return m.Get(key)
}

func Set(key string, val interface{}, ttl time.Duration) {
	m := pool.Get().(*Memory)
	m.Set(key, val, ttl)
	pool.Put(m)
}
