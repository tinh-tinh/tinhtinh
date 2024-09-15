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
	version           string
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
		middlewares: module.Middlewares,
		Dtos:        []Pipe{},
		Security:    []string{},
		module:      module,
		version:     "",
	}
}

func (module *DynamicModule) Use(middleware ...Middleware) *DynamicModule {
	module.Middlewares = append(module.Middlewares, middleware...)

	return module
}

func (module *DynamicModule) Guard(guards ...Guard) *DynamicModule {
	for _, v := range guards {
		mid := module.ParseGuard(v)
		module.Middlewares = append(module.Middlewares, mid)
	}
	return module
}

func (c *DynamicController) Tag(tag string) *DynamicController {
	c.tag = tag
	return c
}

func (c *DynamicController) Version(version string) *DynamicController {
	c.version = version
	return c
}

func (c *DynamicController) Use(middleware ...Middleware) *DynamicController {
	c.middlewares = append(c.middlewares, middleware...)
	return c
}

func (c *DynamicController) Guard(guards ...GuardWithCtrl) *DynamicController {
	for _, v := range guards {
		mid := c.ParseGuardCtrl(v)
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

func (c *DynamicController) Registry() *DynamicController {
	c.globalMiddlewares = append(c.globalMiddlewares, c.middlewares...)
	c.middlewares = []Middleware{}

	return c
}

func (c *DynamicController) Get(path string, handler Handler) {
	c.registry("GET", path, c.ParseCtx(handler))
}

func (c *DynamicController) Post(path string, handler Handler) {
	c.registry("POST", path, c.ParseCtx(handler))
}

func (c *DynamicController) Patch(path string, handler Handler) {
	c.registry("PATCH", path, c.ParseCtx(handler))
}

func (c *DynamicController) Put(path string, handler Handler) {
	c.registry("PUT", path, c.ParseCtx(handler))
}

func (c *DynamicController) Delete(path string, handler Handler) {
	c.registry("DELETE", path, c.ParseCtx(handler))
}

func (c *DynamicController) registry(method string, path string, handler http.Handler) {
	route := ParseRoute(method + " " + path)
	if c.version != "" {
		route.SetPrefix("v" + c.version)
	}
	route.SetPrefix(c.name)

	mergeHandler := handler
	for _, v := range c.middlewares {
		mergeHandler = v(mergeHandler)
	}

	for _, v := range c.globalMiddlewares {
		mergeHandler = v(mergeHandler)
	}

	router := &Router{
		Tag:      c.tag,
		Path:     route.GetPath(),
		Handler:  mergeHandler,
		Dtos:     c.Dtos,
		Security: c.Security,
	}
	c.module.Routers = append(c.module.Routers, router)
	c.free()
}

func (c *DynamicController) free() {
	c.middlewares = []Middleware{}
	c.Dtos = nil
	c.Security = nil
	runtime.GC()
}

func (c *DynamicController) Inject(name Provide) interface{} {
	return c.module.Ref(name)
}
