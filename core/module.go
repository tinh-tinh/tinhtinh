package core

import (
	"net/http"
)

type Mux map[string]http.Handler
type MapValue map[Provide]interface{}

type DocRoute struct {
	Dto      []Pipe
	Security []string
}
type MappingRoute map[string]DocRoute
type MappingDoc map[string]MappingRoute

type DynamicModule struct {
	global      bool
	mux         Mux
	MapperDoc   MappingDoc
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
		MapperDoc:   make(MappingDoc),
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
		for k, v := range mod.mux {
			module.mux[k] = v
		}
		for k, v := range mod.mapperValue {
			module.mapperValue[k] = v
		}
		for k, v := range mod.MapperDoc {
			module.MapperDoc[k] = v
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
		mux:         make(Mux),
		MapperDoc:   make(MappingDoc),
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
		for k, v := range mod.mux {
			newMod.mux[k] = v
		}
		for k, v := range mod.mapperValue {
			newMod.mapperValue[k] = v
		}
		for k, v := range mod.MapperDoc {
			newMod.MapperDoc[k] = v
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
