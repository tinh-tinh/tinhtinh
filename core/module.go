package core

import (
	"net/http"
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

func NewModule(opt NewModuleOptions) *DynamicModule {
	if opt.Scope == "" {
		opt.Scope = Global
	}
	module := &DynamicModule{}
	initModule(module, opt)

	return module
}

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
		module.DataProviders = append(module.DataProviders, mod.getExports()...)
	}

	if module.Scope == Request {
		module.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				module.NewProvider(ProviderOptions{
					Name:  REQUEST,
					Value: r,
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
				h.ServeHTTP(w, r)
				for _, p := range module.getRequest() {
					if p.Value != nil {
						p.Value = nil
					}
				}
			})
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

func (m *DynamicModule) Export(name Provide) *DynamicProvider {
	idx := slices.IndexFunc(m.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == name
	})
	m.DataProviders[idx].Status = PUBLIC
	return m.DataProviders[idx]
}
