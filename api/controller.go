package api

import "net/http"

type Controller struct {
	Name              string
	middlewares       []Middleware
	globalMiddlewares []Middleware
	module            *Module
}

type Handler func(ctx Ctx)

func NewController(name string, module *Module) *Controller {
	return &Controller{
		Name:        name,
		middlewares: []Middleware{},
		module:      module,
	}
}

func (c *Controller) Use(middleware ...Middleware) *Controller {
	c.globalMiddlewares = append(c.globalMiddlewares, middleware...)
	return c
}

func (c *Controller) Guard(guards ...Guard) *Controller {
	for _, v := range guards {
		mid := ParseGuard(v)
		c.middlewares = append(c.middlewares, mid)
	}

	return c
}

func (c *Controller) Pipe(pipe ...Middleware) *Controller {
	c.middlewares = append(c.middlewares, pipe...)
	return c
}

func (c *Controller) Get(path string, handler Handler) {
	c.registry("GET", path, ParseCtx(handler))
}

func (c *Controller) Post(path string, handler Handler) {
	c.registry("POST", path, ParseCtx(handler))
}

func (c *Controller) Patch(path string, handler Handler) {
	c.registry("PATCH", path, ParseCtx(handler))
}

func (c *Controller) Put(path string, handler Handler) {
	c.registry("PUT", path, ParseCtx(handler))
}

func (c *Controller) Delete(path string, handler Handler) {
	c.registry("DELETE", path, ParseCtx(handler))
}

func (c *Controller) registry(method string, path string, handler http.Handler) {
	route := ParseRoute(method + " " + path)
	route.SetPrefix(c.Name)

	mergeHandler := handler
	for _, v := range c.middlewares {
		mergeHandler = v(mergeHandler)
	}

	for _, v := range c.globalMiddlewares {
		mergeHandler = v(mergeHandler)
	}

	c.middlewares = []Middleware{}
	c.module.mux[route.GetPath()] = mergeHandler
}

func (c *Controller) Inject(name string) interface{} {
	return c.module.Ref(name)
}
