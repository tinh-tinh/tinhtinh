package api

import "net/http"

type Module struct {
	middlewares []Middleware
	mux         map[string]http.Handler
	Providers   map[string]interface{}
}

type ModuleParam func() *Module
type ControllerParam func(module *Module) *Controller
type ProviderParam func() Provider

type NewModuleOptions struct {
	Import      []ModuleParam
	Controllers []ControllerParam
	Providers   []ProviderParam
}

func NewModule(opt NewModuleOptions) *Module {
	pd := make(map[string]interface{})

	for _, v := range opt.Providers {
		provider := v()
		pd[provider.Name] = provider.Value
	}

	module := Module{
		middlewares: []Middleware{},
		mux:         make(map[string]http.Handler),
		Providers:   pd,
	}

	for _, v := range opt.Controllers {
		v(&module)
	}

	for _, m := range opt.Import {
		mod := m()
		for k, v := range mod.mux {
			module.mux[k] = v
		}
	}

	return &module
}

func (m *Module) Guard(guard ...Guard) *Module {
	for _, v := range guard {
		mid := ParseGuard(v)
		m.middlewares = append(m.middlewares, mid)
	}

	return m
}

func (m *Module) Interceptor(interceptor ...Middleware) *Module {
	m.middlewares = append(m.middlewares, interceptor...)
	return m
}

func (m *Module) Pipe(pipe ...Middleware) *Module {
	m.middlewares = append(m.middlewares, pipe...)
	return m
}

func (m *Module) Ref(name string) interface{} {
	return m.Providers[name]
}
