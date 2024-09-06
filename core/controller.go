package core

import (
	"net/http"
	"runtime"
	"strings"
)

type Handler func(ctx Ctx)

type DynamicController struct {
	name              string
	tag               string
	Dtos              []Pipe
	Security          []string
	middlewares       []Middleware
	globalMiddlewares []Middleware
	module            *DynamicModule
}

func (module *DynamicModule) NewController(name string) *DynamicController {
	return &DynamicController{
		name:        strings.ToLower(name),
		tag:         name,
		middlewares: []Middleware{},
		Dtos:        []Pipe{},
		Security:    []string{},
		module:      module,
	}
}

func (c *DynamicController) Tag(tag string) *DynamicController {
	c.tag = tag
	return c
}

func (c *DynamicController) Use(middleware ...Middleware) *DynamicController {
	c.globalMiddlewares = append(c.globalMiddlewares, middleware...)
	return c
}

func (c *DynamicController) Guard(guards ...Guard) *DynamicController {
	for _, v := range guards {
		mid := ParseGuard(v)
		c.middlewares = append(c.middlewares, mid)
	}

	return c
}

func (c *DynamicController) GuardWithCtrl(guards ...GuardWithCtrl) *DynamicController {
	for _, v := range guards {
		mid := ParseGuardCtrl(c, v)
		c.middlewares = append(c.middlewares, mid)
	}
	return c
}

func (c *DynamicController) Pipe(dtos ...Pipe) *DynamicController {
	c.Dtos = append(c.Dtos, dtos...)
	middleware := PipeMiddleware(dtos...)
	c.middlewares = append(c.middlewares, middleware)
	return c
}

func (c *DynamicController) AddSecurity(security ...string) *DynamicController {
	c.Security = append(c.Security, security...)
	return c
}

func (c *DynamicController) Get(path string, handler Handler) {
	c.registry("GET", path, ParseCtx(handler))
}

func (c *DynamicController) Post(path string, handler Handler) {
	c.registry("POST", path, ParseCtx(handler))
}

func (c *DynamicController) Patch(path string, handler Handler) {
	c.registry("PATCH", path, ParseCtx(handler))
}

func (c *DynamicController) Put(path string, handler Handler) {
	c.registry("PUT", path, ParseCtx(handler))
}

func (c *DynamicController) Delete(path string, handler Handler) {
	c.registry("DELETE", path, ParseCtx(handler))
}

func (c *DynamicController) registry(method string, path string, handler http.Handler) {
	route := ParseRoute(method + " " + path)
	route.SetPrefix(c.name)

	mergeHandler := handler
	for _, v := range c.middlewares {
		mergeHandler = v(mergeHandler)
	}

	for _, v := range c.globalMiddlewares {
		mergeHandler = v(mergeHandler)
	}

	c.module.mux[route.GetPath()] = mergeHandler
	if c.module.MapperDoc[c.tag] == nil {
		c.module.MapperDoc[c.tag] = make(map[string]DocRoute)
	}

	ct := c.module.MapperDoc[c.tag]
	docRoute := DocRoute{
		Dto:      c.Dtos,
		Security: c.Security,
	}
	ct[route.GetPath()] = docRoute

	c.middlewares = []Middleware{}
	c.Dtos = nil
	c.Security = nil
	runtime.GC()
}

func (c *DynamicController) Inject(name Provide) interface{} {
	return c.module.Ref(name)
}
