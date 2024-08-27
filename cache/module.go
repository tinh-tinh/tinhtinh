package cache

type CacheType string

const (
	InMemory CacheType = "memory"
	Redis    CacheType = "redis"
)

type ModuleOptions struct {
	Type CacheType
}

type Module struct {
	memory *Memory
}

func Register(opt ModuleOptions) *Module {
	if opt.Type == InMemory {
		return &Module{
			memory: NewInMemory(),
		}
	}
	return nil
}
