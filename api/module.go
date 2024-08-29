package api

import "net/http"

type Module struct {
	middlewares []Middleware
	mux         map[string]http.Handler
}

type Factory func(params ...interface{})

func NewModule() *Module {
	return &Module{
		middlewares: []Middleware{},
		mux:         make(map[string]http.Handler),
	}
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

func (m *Module) Import(modules ...*Module) *Module {
	for _, mo := range modules {
		for k, v := range mo.mux {
			var mergeHandler = v
			for _, v := range m.middlewares {
				mergeHandler = v(mergeHandler)
			}
			m.mux[k] = mergeHandler
		}
		mo = nil
	}
	return m
}

func (m *Module) Controllers(controllers ...*Controller) *Module {
	for _, c := range controllers {
		for k, v := range c.mux {
			var mergeHandler = v
			for _, v := range m.middlewares {
				mergeHandler = v(mergeHandler)
			}
			m.mux[k] = mergeHandler
		}
		c = nil
	}
	return m
}
