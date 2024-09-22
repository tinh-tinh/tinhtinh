package core

type Provide string

const REQUEST Provide = "REQUEST"

type ProvideStatus string

const (
	PUBLIC  ProvideStatus = "public"
	PRIVATE ProvideStatus = "private"
)

type Factory func(param ...interface{}) interface{}

type DynamicProvider struct {
	Scope   Scope
	Status  ProvideStatus
	Name    Provide
	Value   interface{}
	factory Factory
	inject  []Provide
}

type ProviderOptions struct {
	Name    Provide
	Value   interface{}
	Factory Factory
	Inject  []Provide
}

func (module *DynamicModule) NewProvider(opt ProviderOptions) *DynamicProvider {
	provider := &DynamicProvider{
		Name:   opt.Name,
		Status: PRIVATE,
		Scope:  module.Scope,
	}
	module.DataProviders = append(module.DataProviders, provider)
	if provider.Scope == Request {
		provider.inject = opt.Inject
		provider.factory = opt.Factory
		provider.Value = opt.Value
		return provider
	}
	provider.Value = opt.Value
	if opt.Value == nil {
		var values []interface{}
		for _, p := range opt.Inject {
			values = append(values, module.ref(p))
		}
		provider.Value = opt.Factory(values...)
	}

	return provider
}

func (module *DynamicModule) getExports() []*DynamicProvider {
	exports := make([]*DynamicProvider, 0)
	for _, v := range module.DataProviders {
		if v.Status == PUBLIC {
			exports = append(exports, v)
		}
	}

	return exports
}

func (module *DynamicModule) getRequest() []*DynamicProvider {
	reqs := make([]*DynamicProvider, 0)
	for _, v := range module.DataProviders {
		if v.Scope == Request {
			reqs = append(reqs, v)
		}
	}
	return reqs
}
