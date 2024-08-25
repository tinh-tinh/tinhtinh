package tinhtinh

import "net/http"

type Controller struct {
	Name              string
	middlewares       []Middleware
	globalMiddlewares []Middleware
	mux               map[string]http.Handler
}

type Handler func(ctx Ctx)

func NewController(name string) *Controller {
	return &Controller{
		Name:        name,
		middlewares: []Middleware{},
		mux:         make(map[string]http.Handler),
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
	c.mux[route.GetPath()] = mergeHandler
}
