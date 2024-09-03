package core

import (
	"net/http"
)

type Mux map[string]http.Handler
type MapValue map[Provide]interface{}

type DynamicModule struct {
	global      bool
	MapMux      map[string]Mux
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
		MapMux:      make(map[string]Mux),
		mapperValue: make(MapValue),
		global:      opt.Global,
	}

	// Providers
	providers := make([]Provider, 0)
	for _, p := range opt.Providers {
		p(module)
		providers = append(providers, p)
	}

	// Imports
	for _, m := range opt.Imports {
		mod := m(module)
		for k, v := range mod.MapMux {
			module.MapMux[k] = v
		}
		for k, v := range mod.mapperValue {
			module.mapperValue[k] = v
		}

		if module.global {
			mod.Providers(providers...)
		}
	}

	// Controllers
	for _, ct := range opt.Controllers {
		ct(module)
	}

	return module
}

func (m *DynamicModule) New(opt NewModuleOptions) *DynamicModule {
	newMod := &DynamicModule{
		mapperValue: m.mapperValue,
		MapMux:      make(map[string]Mux),
		global:      opt.Global,
	}

	// Providers
	providers := make([]Provider, 0)
	for _, p := range opt.Providers {
		p(newMod)
		providers = append(providers, p)
	}

	// Imports
	for _, m := range opt.Imports {
		mod := m(newMod)
		for k, v := range mod.MapMux {
			newMod.MapMux[k] = v
		}
		for k, v := range mod.mapperValue {
			newMod.mapperValue[k] = v
		}

		if newMod.global {
			mod.Providers(providers...)
		}
	}

	// Controllers
	for _, ct := range opt.Controllers {
		ct(newMod)
	}

	return newMod
}

func (m *DynamicModule) Controllers(controllers ...Controller) *DynamicModule {
	for _, v := range controllers {
		v(m)
	}
	return m
}

func (m *DynamicModule) Providers(providers ...Provider) {
	for _, v := range providers {
		v(m)
	}
}

func (m *DynamicModule) Ref(name Provide) interface{} {
	return m.mapperValue[name]
}
