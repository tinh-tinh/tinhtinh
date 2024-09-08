package core

import (
	"errors"
	"net/http"
	"runtime"
	"slices"
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
	Middlewares []Middleware
	MapperDoc   MappingDoc
	providers   []*DynamicProvider
	Exports     []*DynamicProvider
}

type Module func(module *DynamicModule) *DynamicModule
type Controller func(module *DynamicModule) *DynamicController
type Provider func(module *DynamicModule) *DynamicProvider

type NewModuleOptions struct {
	Global      bool
	Imports     []Module
	Controllers []Controller
	Providers   []Provider
	Exports     []Provider
}

func NewModule(opt NewModuleOptions) *DynamicModule {
	module := &DynamicModule{
		mux:       make(Mux),
		MapperDoc: make(MappingDoc),
		providers: []*DynamicProvider{},
		Exports:   []*DynamicProvider{},
		global:    opt.Global,
	}

	// Providers
	for _, p := range opt.Providers {
		p(module)
	}

	// Imports
	for _, m := range opt.Imports {
		mod := m(module)
		for k, v := range mod.mux {
			module.mux[k] = v
		}
		module.providers = append(module.providers, mod.Exports...)
		module.Exports = append(module.providers, mod.Exports...)
		for k, v := range mod.MapperDoc {
			module.MapperDoc[k] = v
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
		providers: []*DynamicProvider{},
		Exports:   []*DynamicProvider{},
		mux:       make(Mux),
		MapperDoc: make(MappingDoc),
		global:    opt.Global,
	}

	newMod.providers = append(newMod.providers, m.Exports...)
	// Providers
	providers := make([]Provider, 0)
	for _, p := range opt.Providers {
		p(newMod)
		providers = append(providers, p)
	}

	// Imports
	for _, mFnc := range opt.Imports {
		mod := mFnc(newMod)
		for k, v := range mod.mux {
			newMod.mux[k] = v
		}
		newMod.providers = append(newMod.providers, mod.Exports...)
		newMod.Exports = append(newMod.providers, mod.Exports...)
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

	runtime.GC()
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
	idx := slices.IndexFunc(m.providers, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	return m.providers[idx].Value
}

func (m *DynamicModule) Export(key Provide) {
	idx := slices.IndexFunc(m.providers, func(e *DynamicProvider) bool {
		return e.Name == key
	})
	if idx == -1 {
		panic(errors.New("key of provider not found"))
	}
	m.Exports = append(m.Exports, m.providers[idx])
}
