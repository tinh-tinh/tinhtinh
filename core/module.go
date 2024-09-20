package core

import (
	"errors"
	"runtime"
	"slices"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type DocRoute struct {
	Dto      []Pipe
	Security []string
}

type DynamicModule struct {
	global      bool
	Routers     []*Router
	Middlewares []Middleware
	providers   []*DynamicProvider
	Exports     []*DynamicProvider
	hooks       []HookModule
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
		utils.Log(
			utils.Green("[TT] "),
			utils.White(time.Now().Format("2006-01-02 15:04:05")),
			utils.Yellow(" [Module Initializer] "),
			utils.Green(utils.GetFunctionName(m)+"\n"),
		)

		mod.init()
		module.Routers = append(module.Routers, mod.Routers...)
		module.providers = append(module.providers, mod.Exports...)
		module.Exports = append(module.providers, mod.Exports...)
	}

	// Controllers
	for _, ct := range opt.Controllers {
		ct(module)
	}

	return module
}

func (m *DynamicModule) New(opt NewModuleOptions) *DynamicModule {
	newMod := &DynamicModule{
		global: opt.Global,
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
		utils.Log(
			utils.Green("[TT] "),
			utils.White(time.Now().Format("2006-01-02 15:04:05")),
			utils.Yellow(" [Module Initializer] "),
			utils.Green(utils.GetFunctionName(m)+"\n"),
		)

		mod.init()
		newMod.Routers = append(newMod.Routers, mod.Routers...)
		newMod.providers = append(newMod.providers, mod.Exports...)
		newMod.Exports = append(newMod.providers, mod.Exports...)

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

func (m *DynamicModule) RefFactory(name Provide, ctx Ctx) interface{} {
	idx := slices.IndexFunc(m.providers, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	return m.providers[idx].Factory(ctx)
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
