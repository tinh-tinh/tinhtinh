package core

import (
	"runtime"
	"slices"
)

type DocRoute struct {
	Dto      []Pipe
	Security []string
}

type DynamicModule struct {
	global        bool
	Routers       []*Router
	Middlewares   []Middleware
	DataProviders []*DynamicProvider
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
		global: opt.Global,
	}

	// Providers
	for _, p := range opt.Providers {
		p(module)
	}

	// Imports
	for _, m := range opt.Imports {
		mod := m(module)
		module.Routers = append(module.Routers, mod.Routers...)
		module.DataProviders = append(module.DataProviders, mod.getExports()...)
	}

	// Controllers
	for _, ct := range opt.Controllers {
		ct(module)
	}

	// Exports
	for _, e := range opt.Exports {
		provider := e(module)
		provider.Status = PUBLIC
	}

	return module
}

func (m *DynamicModule) New(opt NewModuleOptions) *DynamicModule {
	newMod := &DynamicModule{
		global: opt.Global,
	}

	newMod.DataProviders = append(newMod.DataProviders, m.getExports()...)
	// Providers
	providers := make([]Provider, 0)
	for _, p := range opt.Providers {
		p(newMod)
		providers = append(providers, p)
	}

	// Imports
	for _, mFnc := range opt.Imports {
		mod := mFnc(newMod)
		newMod.Routers = append(newMod.Routers, mod.Routers...)
		newMod.DataProviders = append(newMod.DataProviders, mod.getExports()...)

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
	idx := slices.IndexFunc(m.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	return m.DataProviders[idx].Value
}

func (m *DynamicModule) Export(name Provide) *DynamicProvider {
	idx := slices.IndexFunc(m.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	m.DataProviders[idx].Status = PUBLIC
	return m.DataProviders[idx]
}
