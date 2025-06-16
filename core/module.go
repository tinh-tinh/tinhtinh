package core

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/common/color"
)

type DocRoute struct {
	Dto      []PipeDto
	Security []string
}

type Scope string

const (
	Global    Scope = "global"
	Request   Scope = "request"
	Transient Scope = "transient"
)

type Module interface {
	OnInit(hooks ...HookModule) Module
	New(opt NewModuleOptions) Module
	Controllers(controllers ...Controllers) Module
	Providers(providers ...Providers) Module
	Export(name Provide) Provider
	Ref(name Provide, ctx ...Ctx) interface{}
	findIdx(name Provide) int
	init()
	GetRouters() []*Router
	GetExports() []Provider
	free()
	NewController(name string) Controller
	NewProvider(opt ProviderParams) Provider
	Consumer(consumer *Consumer) Module
	Guard(guards ...Guard) Module
	Use(middleware ...Middleware) Module
	UseRef(middlewareRefs ...MiddlewareRef) Module
	GetDataProviders() []Provider
	AppendDataProviders(providers ...Provider)
	GetScope() Scope
}

type DynamicModule struct {
	sync.Pool
	isRoot        bool
	Scope         Scope
	Routers       []*Router
	Middlewares   []Middleware
	DataProviders []Provider
	hooks         []HookModule
	interceptor   Interceptor
}

type Modules func(module Module) Module
type Controllers func(module Module) Controller
type Providers func(module Module) Provider

type NewModuleOptions struct {
	Scope       Scope
	Imports     []Modules
	Controllers []Controllers
	Providers   []Providers
	Exports     []Providers
	Guards      []Guard
	Middlewares []Middleware
	Interceptor Interceptor
}

// NewModule creates a new module with the given options.
//
// The scope of the module will default to Global if not specified.
// The module will be initialized with the given imports, controllers, providers and exports.
func NewModule(opt NewModuleOptions) Module {
	if opt.Scope == "" {
		opt.Scope = Global
	}
	module := &DynamicModule{isRoot: true}
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
func (m *DynamicModule) New(opt NewModuleOptions) Module {
	if opt.Scope == "" {
		opt.Scope = Global
	}
	newMod := &DynamicModule{isRoot: false}
	newMod.DataProviders = append(newMod.DataProviders, m.GetExports()...)
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
	// Parse middleware
	module.Middlewares = append(module.Middlewares, opt.Middlewares...)

	// Parse guards
	for _, g := range opt.Guards {
		if g == nil {
			continue
		}
		mid := module.ParseGuard(g)
		module.Middlewares = append(module.Middlewares, mid)
	}

	// Parse interceptor
	module.interceptor = opt.Interceptor

	// Imports
	for _, m := range opt.Imports {
		if m == nil {
			continue
		}
		mod := m(module)
		fmt.Printf("%s %s %s %s\n",
			color.Green("[TT]"),
			color.White(time.Now().Format("2006-01-02 15:04:05")),
			color.Yellow("[Module Initializer]"),
			color.Green(common.GetFunctionName(m)),
		)

		mod.init()
		module.Routers = append(module.Routers, mod.GetRouters()...)
		module.appendProvider(mod.GetExports()...)
		mod.free()
	}

	// Providers
	for _, p := range opt.Providers {
		if p == nil {
			continue
		}
		p(module)
	}

	isRequest := slices.ContainsFunc(module.DataProviders, func(e Provider) bool {
		return e.GetScope() == Request
	})

	if module.Scope == Request || isRequest {
		module.Use(requestMiddleware(module))
	}

	// Controllers
	for _, ct := range opt.Controllers {
		if ct == nil {
			continue
		}
		ct(module)
	}

	// Exports
	for _, e := range opt.Exports {
		if e == nil {
			continue
		}
		provider := e(module)
		provider.SetStatus(PUBLIC)
	}
}

// Controllers registers the given controllers with the module.
// The controllers are registered in the order they are given.
func (m *DynamicModule) Controllers(controllers ...Controllers) Module {
	for _, v := range controllers {
		v(m)
	}
	return m
}

// Providers registers the given providers with the module.
// The providers are registered in the order they are given.
func (m *DynamicModule) Providers(providers ...Providers) Module {
	for _, v := range providers {
		v(m)
	}

	return m
}

// Ref returns the value of the provider with the given name.
// If the provider is not found, Ref returns nil.
func (m *DynamicModule) Ref(name Provide, ctx ...Ctx) interface{} {
	if name == REQUEST {
		return ctx[0].Req()
	}
	idx := slices.IndexFunc(m.DataProviders, func(e Provider) bool {
		return e.GetName() == name
	})
	if idx == -1 {
		return nil
	}

	prd := m.DataProviders[idx]
	if prd.GetScope() == Request {
		if len(ctx) == 0 {
			panic("request provider need ctx as parameters")
		}
		return ctx[0].Get(name)
	} else if prd.GetScope() == Transient {
		var values []interface{}
		for _, p := range prd.GetInject() {
			values = append(values, m.Ref(p))
		}
		return prd.GetFactory()(values...)
	}
	return prd.GetValue()
}

func (m *DynamicModule) findIdx(name Provide) int {
	idx := slices.IndexFunc(m.DataProviders, func(e Provider) bool {
		return e.GetName() == name
	})
	return idx
}

// Export sets the status of the provider with the given name to PUBLIC and returns
// the provider.
func (m *DynamicModule) Export(name Provide) Provider {
	idx := slices.IndexFunc(m.DataProviders, func(e Provider) bool {
		return e.GetName() == name
	})
	m.DataProviders[idx].SetStatus(PUBLIC)
	return m.DataProviders[idx]
}

func (m *DynamicModule) GetRouters() []*Router {
	return m.Routers
}

func (m *DynamicModule) free() {
	m.Routers = nil
}

func requestMiddleware(module *DynamicModule) Middleware {
	return func(ctx Ctx) error {
		for _, p := range module.getRequest() {
			if p.GetValue() == nil {
				var values []interface{}
				for _, p := range p.GetInject() {
					values = append(values, module.Ref(p, ctx))
				}

				factory := p.GetFactory()
				value := factory(values...)
				ctx.Set(p.GetName(), value)
			}
		}
		return ctx.Next()
	}
}

func (m *DynamicModule) GetDataProviders() []Provider {
	return m.DataProviders
}

func (m *DynamicModule) AppendDataProviders(providers ...Provider) {
	m.DataProviders = append(m.DataProviders, providers...)
}

func (m *DynamicModule) GetScope() Scope {
	return m.Scope
}
