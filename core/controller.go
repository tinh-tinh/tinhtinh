package core

import (
	"net/http"
	"strings"
)

type Handler func(ctx Ctx)

type DynamicController struct {
	name              string
	tag               string
	middlewares       []Middleware
	globalMiddlewares []Middleware
	module            *DynamicModule
}

func NewController(name string, module *DynamicModule) *DynamicController {
	return &DynamicController{
		name:        strings.ToLower(name),
		tag:         name,
		middlewares: []Middleware{},
		module:      module,
	}
}

func (c *DynamicController) Tag(tag string) {
	c.tag = tag
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

func (c *DynamicController) Pipe(pipe ...Middleware) *DynamicController {
	c.middlewares = append(c.middlewares, pipe...)
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

	c.middlewares = []Middleware{}
	if c.module.MapMux[c.tag] == nil {
		c.module.MapMux[c.tag] = make(Mux)
	}
	mux := c.module.MapMux[c.tag]
	mux[route.GetPath()] = mergeHandler
	c.module.MapMux[c.tag] = mux
}

func (c *DynamicController) Inject(name Provide) interface{} {
	return c.module.Ref(name)
}
