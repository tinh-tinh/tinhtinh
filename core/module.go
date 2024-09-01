package core

import (
	"net/http"
)

type Mux map[string]http.Handler
type MapValue map[Provide]interface{}

type DynamicModule struct {
	global      bool
	mux         Mux
	mapperValue MapValue
}

type Module func(module *DynamicModule) *DynamicModule
type Controller func(module *DynamicModule) *DynamicController
type Provider func(module *DynamicModule) *DynamicProvider

type NewModuleOptions struct {
	Global      bool
	Imports     []Module
	Controllers []Controller
	Providers   []Provider
}

func NewModule(opt NewModuleOptions) *DynamicModule {
	module := &DynamicModule{
		mux:         make(Mux),
		mapperValue: make(MapValue),
		global:      opt.Global,
	}

	// Providers
	providers := make([]*DynamicProvider, 0)
	for _, p := range opt.Providers {
		provider := p(module)
		providers = append(providers, provider)
	}
	module.setProviders(providers...)

	// Imports
	for _, m := range opt.Imports {
		mod := m(module)
		for k, v := range mod.mux {
			module.mux[k] = v
		}
		for k, v := range mod.mapperValue {
			module.mapperValue[k] = v
		}

		if module.global {
			mod.setProviders(providers...)
		}
	}

	// Controllers
	for _, ct := range opt.Controllers {
		ct(module)
	}

	return module
}

func (m *DynamicModule) setProviders(providers ...*DynamicProvider) {
	for _, v := range providers {
		m.mapperValue[v.Name] = v.Value
	}
}

func (m *DynamicModule) Ref(name Provide) interface{} {
	return m.mapperValue[name]
}
