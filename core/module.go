package core

import (
	"slices"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type DocRoute struct {
	Dto      []Pipe
	Security []string
}

type Scope string

const (
	Global  Scope = "global"
	Request Scope = "request"
)

type DynamicModule struct {
	Scope         Scope
	Routers       []*Router
	Middlewares   []Middleware
	DataProviders []*DynamicProvider
	hooks         []HookModule
}

type Module func(module *DynamicModule) *DynamicModule
type Controller func(module *DynamicModule) *DynamicController
type Provider func(module *DynamicModule) *DynamicProvider

type NewModuleOptions struct {
	Scope       Scope
	Imports     []Module
	Controllers []Controller
	Providers   []Provider
	Exports     []Provider
}

// NewModule creates a new module with the given options.
//
// The scope of the module will default to Global if not specified.
// The module will be initialized with the given imports, controllers, providers and exports.
func NewModule(opt NewModuleOptions) *DynamicModule {
	if opt.Scope == "" {
		opt.Scope = Global
	}
	module := &DynamicModule{}
	initModule(module, opt)

	return module
}

// New creates a new module as a sub-module of the current module.
//
// The sub-module will inherit all the exports of the current module and
// will have the same middlewares as the current module.
//
// The scope of the sub-module will default to Global if not specified.
// The sub-module will be initialized with the given imports, controllers, providers and exports.
func (m *DynamicModule) New(opt NewModuleOptions) *DynamicModule {
	if opt.Scope == "" {
		opt.Scope = Global
	}
	newMod := &DynamicModule{}
	newMod.DataProviders = append(newMod.DataProviders, m.getExports()...)
	newMod.Middlewares = append(newMod.Middlewares, m.Middlewares...)

	initModule(newMod, opt)
	return newMod
}

// initModule initializes the given module with the given options.
//
// It sets the scope of the module, runs the providers, imports the sub-modules,
// runs the controllers, and sets the exports.
//
// If the scope of the module is Request, it wraps the handler with a middleware
// that creates a new request provider for each request, and sets the value of
// the providers that are injected with the request to the value of the provider
// of the current request. After the request is handled, it sets the value of
// the providers that are injected with the request to nil.
func initModule(module *DynamicModule, opt NewModuleOptions) {
	module.Scope = opt.Scope
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
		module.appendProvider(mod.getExports()...)
	}

	if module.Scope == Request {
		module.Use(func(ctx Ctx) error {
			module.NewProvider(ProviderOptions{
				Name:  REQUEST,
				Value: ctx.Req(),
			})
			for _, p := range module.getRequest() {
				if p.Value == nil {
					var values []interface{}
					for _, p := range p.inject {
						values = append(values, module.Ref(p))
					}

					p.Value = p.factory(values...)
				}
			}
			err := ctx.Next()
			if err != nil {
				return err
			}
			for _, p := range module.getRequest() {
				if p.Value != nil {
					p.Value = nil
				}
			}
			return nil
		})
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
}

// Controllers registers the given controllers with the module.
// The controllers are registered in the order they are given.
func (m *DynamicModule) Controllers(controllers ...Controller) *DynamicModule {
	for _, v := range controllers {
		v(m)
	}
	return m
}

// Providers registers the given providers with the module.
// The providers are registered in the order they are given.
func (m *DynamicModule) Providers(providers ...Provider) {
	for _, v := range providers {
		v(m)
	}
}

// Ref returns the value of the provider with the given name.
// If the provider is not found, Ref returns nil.
func (m *DynamicModule) Ref(name Provide) interface{} {
	idx := slices.IndexFunc(m.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	if idx == -1 {
		return nil
	}
	return m.DataProviders[idx].Value
}

func (m *DynamicModule) findIdx(name Provide) int {
	idx := slices.IndexFunc(m.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	return idx
}

// Export sets the status of the provider with the given name to PUBLIC and returns
// the provider.
func (m *DynamicModule) Export(name Provide) *DynamicProvider {
	idx := slices.IndexFunc(m.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	m.DataProviders[idx].Status = PUBLIC
	return m.DataProviders[idx]
}
