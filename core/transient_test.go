package core_test

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type Container struct {
	pools map[reflect.Type]*sync.Pool
	mu    sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		pools: make(map[reflect.Type]*sync.Pool),
	}
}

func RegisterTransient[T any](c *Container, factory func() T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	typ := reflect.TypeOf((*T)(nil)).Elem()
	c.pools[typ] = &sync.Pool{
		New: func() any {
			return factory()
		},
	}
}

func Resolve[T any](c *Container) T {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	pool := c.pools[typ]
	if pool == nil {
		panic("type not registered")
	}
	return pool.Get().(T)
}

func Release[T any](c *Container, obj T) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	pool := c.pools[typ]
	if pool == nil {
		panic("type not registered")
	}
	pool.Put(obj)
}

type UserService struct {
	Name string
}

func Test_Transient(t *testing.T) {
	c := NewContainer()

	RegisterTransient(c, func() *UserService {
		return &UserService{}
	})

	user := Resolve[*UserService](c)
	user.Name = "Alice"
	fmt.Println(user)

	Release(c, user)
}
