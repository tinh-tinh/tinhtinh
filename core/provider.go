package core

import (
	"slices"

	"github.com/tinh-tinh/tinhtinh/v2/common"
)

type Provide string

const REQUEST Provide = "REQUEST"

type ProvideStatus string

const (
	PUBLIC  ProvideStatus = "public"
	PRIVATE ProvideStatus = "private"
)

type Factory func(param ...interface{}) interface{}

type ProviderType string

type Provider interface {
	GetName() Provide
	SetName(name Provide)
	GetValue() interface{}
	SetValue(value interface{})
	GetStatus() ProvideStatus
	SetStatus(status ProvideStatus)
	GetScope() Scope
	SetScope(scope Scope)
	SetInject(inject []Provide)
	GetInject() []Provide
	SetFactory(factory Factory)
	GetFactory() Factory
}

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

func (p *DynamicProvider) GetName() Provide {
	return p.Name
}

func (p *DynamicProvider) SetName(name Provide) {
	p.Name = name
}

func (p *DynamicProvider) GetValue() interface{} {
	return p.Value
}

func (p *DynamicProvider) SetValue(value interface{}) {
	p.Value = value
}

func (p *DynamicProvider) GetStatus() ProvideStatus {
	return p.Status
}

func (p *DynamicProvider) SetStatus(status ProvideStatus) {
	p.Status = status
}

func (p *DynamicProvider) GetScope() Scope {
	return p.Scope
}

func (p *DynamicProvider) SetScope(scope Scope) {
	p.Scope = scope
}

func (p *DynamicProvider) SetInject(inject []Provide) {
	p.inject = inject
}

func (p *DynamicProvider) GetInject() []Provide {
	return p.inject
}

func (p *DynamicProvider) SetFactory(factory Factory) {
	p.factory = factory
}

func (p *DynamicProvider) GetFactory() Factory {
	return p.factory
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
	Type   ProviderType
}

type ProviderParams interface {
	ProviderOptions | interface{}
}

func (module *DynamicModule) NewProvider(param ProviderParams) Provider {
	providerOptions, ok := param.(ProviderOptions)
	if ok {
		return InitProviders(module, providerOptions)
	}
	nameProvide := common.GetStructName(param)
	options := ProviderOptions{
		Name:  Provide(nameProvide),
		Value: param,
	}
	return InitProviders(module, options)
}

// getExports returns a list of providers that are exported by the module.
// The exported providers are the providers that have the status PUBLIC.
func (module *DynamicModule) GetExports() []Provider {
	exports := make([]Provider, 0)
	for _, v := range module.DataProviders {
		if v.GetStatus() == PUBLIC {
			exports = append(exports, v)
		}
	}

	return exports
}

// getRequest returns a list of providers that have the scope Request.
// The providers are the providers that will be injected with the request.
func (module *DynamicModule) getRequest() []Provider {
	reqs := make([]Provider, 0)
	for _, v := range module.DataProviders {
		if v.GetScope() == Request {
			reqs = append(reqs, v)
		}
	}
	return reqs
}

// appendProvider appends the given providers to the module's list of providers.
// If the provider already exists with the same name, it will override the existing
// provider.
func (module *DynamicModule) appendProvider(providers ...Provider) {
	for _, provider := range providers {
		idx := module.findIdx(provider.GetName())
		if idx == -1 {
			module.DataProviders = append(module.DataProviders, provider)
			continue
		}
		module.DataProviders[idx] = provider
	}
}

func InitProviders(module Module, opt ProviderOptions) Provider {
	var provider Provider
	providerIdx := module.findIdx(opt.Name)
	if providerIdx != -1 {
		provider = module.GetDataProviders()[providerIdx]
	} else {
		provider = &DynamicProvider{
			Name:   opt.Name,
			Status: PRIVATE,
			Scope:  opt.Scope,
		}
		module.AppendDataProviders(provider)
	}

	if provider.GetScope() == "" {
		provider.SetScope(module.GetScope())
	}

	// Handle transient
	if provider.GetScope() == Transient {
		provider.SetInject(opt.Inject)
		provider.SetFactory(opt.Factory)
		return provider
	}

	// Handle request scope
	reqInject := slices.ContainsFunc(opt.Inject, func(p Provide) bool {
		return p == REQUEST
	})
	if reqInject {
		provider.SetScope(Request)
	}

	if provider.GetScope() == Request {
		provider.SetInject(opt.Inject)
		provider.SetFactory(opt.Factory)
		provider.SetValue(opt.Value)
		return provider
	}

	// Handle singleton
	provider.SetValue(opt.Value)
	if opt.Value == nil {
		var values []interface{}
		for _, p := range opt.Inject {
			values = append(values, module.Ref(p))
		}
		val := opt.Factory(values...)
		if val != nil {
			provider.SetValue(val)
		} else {
			provider.SetFactory(opt.Factory)
		}
	}

	return provider
}

func Inject[P any](module RefProvider) *P {
	var provider P
	name := common.GetStructName(provider)
	svc, ok := module.Ref(Provide(name)).(*P)
	if !ok {
		return nil
	}
	return svc
}
