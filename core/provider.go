package core

import "slices"

type Provider interface {
	GetName() Provide
	GetValue() interface{}
	IsPublic() bool
	IsRequest() bool
}

type Provide string

const REQUEST Provide = "REQUEST"

type ProvideStatus string

const (
	PUBLIC  ProvideStatus = "public"
	PRIVATE ProvideStatus = "private"
)

type Factory func(param ...interface{}) interface{}

type DynamicProvider struct {
	// Scope of the provider. Default is Global.
	Scope Scope
	// Status of the provider. Default is PRIVATE.
	Status ProvideStatus
	// Name of the provider.
	Name Provide
	// Value of the provider.
	Value interface{}
	// Factory function for retrieving the value of the other providers in the module.
	factory Factory
	// Providers that are injected with the provider.
	inject []Provide
}

type ProviderOptions struct {
	// Scope of the provider. Default is Global.
	Scope Scope
	// Name of the provider.
	Name Provide
	// Value of the provider.
	Value interface{}
	// Factory function for retrieving the value of the other providers in the module.
	// If the factory function is nil, the value of the provider will be set to the
	// given value.
	Factory Factory
	// Providers that are injected with the provider.
	Inject []Provide
}

// NewProvider creates a new provider with the given options.
// If the provider with the same name has existed, the value of the provider
// will be override.
// If the scope of the provider is Request, the provider will be injected with
// the given injects and the value of the provider will be set to the result of
// the factory function with the given injects.
// Otherwise, the value of the provider will be set to the given value, or the
// result of the factory function with the given injects if the value is nil.
func (module *DynamicModule) NewProvider(opt ProviderOptions) *DynamicProvider {
	var provider *DynamicProvider
	providerIdx := module.findIdx(opt.Name)
	if providerIdx != -1 {
		provider = module.DataProviders[providerIdx]
	} else {
		provider = &DynamicProvider{
			Name:   opt.Name,
			Status: PRIVATE,
			Scope:  opt.Scope,
		}
		module.DataProviders = append(module.DataProviders, provider)
	}
	reqInject := slices.ContainsFunc(opt.Inject, func(p Provide) bool {
		return p == REQUEST
	})
	if reqInject {
		provider.Scope = Request
	}
	if provider.Scope == "" {
		provider.Scope = module.Scope
	}
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
			values = append(values, module.Ref(p))
		}
		provider.Value = opt.Factory(values...)
	}

	return provider
}

// getExports returns a list of providers that are exported by the module.
// The exported providers are the providers that have the status PUBLIC.
func (module *DynamicModule) getExports() []*DynamicProvider {
	exports := make([]*DynamicProvider, 0)
	for _, v := range module.DataProviders {
		if v.Status == PUBLIC {
			exports = append(exports, v)
		}
	}

	return exports
}

// getRequest returns a list of providers that have the scope Request.
// The providers are the providers that will be injected with the request.
func (module *DynamicModule) getRequest() []*DynamicProvider {
	reqs := make([]*DynamicProvider, 0)
	for _, v := range module.DataProviders {
		if v.Scope == Request {
			reqs = append(reqs, v)
		}
	}
	return reqs
}

// appendProvider appends the given providers to the module's list of providers.
// If the provider already exists with the same name, it will override the existing
// provider.
func (module *DynamicModule) appendProvider(providers ...*DynamicProvider) {
	for _, provider := range providers {
		idx := module.findIdx(provider.Name)
		if idx == -1 {
			module.DataProviders = append(module.DataProviders, provider)
			continue
		}
		module.DataProviders[idx] = provider
	}
}
