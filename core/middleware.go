package core

import (
	"fmt"
	"net/http"
)

type Middleware func(ctx Ctx) error
type MiddlewareRef func(ref RefProvider, ctx Ctx) error

type RefProvider interface {
	Ref(name Provide, ctx ...Ctx) interface{}
}

type middlewareRaw func(http.Handler) http.Handler

// ParseCtxMiddleware wraps a Middleware function and returns a middlewareRaw
// that can be used by http server. It provides a Ctx instance to the wrapped
// middleware function and automatically sets the handler of the Ctx instance.
func ParseCtxMiddleware(app *App, ctxMid Middleware) middlewareRaw {
	ctx := app.pool.Get().(*Ctx)
	defer app.pool.Put(ctx)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.SetCtx(w, r)
			ctx.SetHandler(h)
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
					app.errorHandler(err, *ctx)
				}
			}()
			err = ctxMid(*ctx)
			if err != nil {
				app.errorHandler(err, *ctx)
				return
			}
		})
	}
}

// ParseMiddlewareRef wraps a MiddlewareRef function and returns a Middleware
// that can be used by the controller. The wrapped MiddlewareRef function is
// called with the controller as the first argument and the Ctx as the second
// argument. The middleware handler returned by the wrapped MiddlewareRef
// function is run by the controller's Use method. The middleware handler is run
// before the controller's handlers. The controller's middleware handlers are run
// after the module's middleware handlers. The module middleware handlers are run
// after the module's parent middleware handlers. The module middleware handlers
// are run before the module's controllers. The controller's middleware handlers
// are run before the controller's handlers. If the middleware handler returns an
// error, the request is rejected with the error.
func (c *DynamicController) ParseMiddlewareRef(ref MiddlewareRef) Middleware {
	return func(ctx Ctx) error {
		return ref(c, ctx)
	}
}

// Use appends the given middleware functions to the controller's list of
// middleware handlers. The middleware handlers are run in the order they
// are added to the controller. The middleware handlers are run before the
// controller's handlers. The controller's middleware handlers are run
// after the module's middleware handlers. The module middleware handlers
// are run after the module's parent middleware handlers. The module
// middleware handlers are run before the module's controllers. The
// controller middleware handlers are run before the controller's
// handlers.
func (c *DynamicController) Use(middleware ...Middleware) *DynamicController {
	c.middlewares = append(c.middlewares, middleware...)
	return c
}

// UseRef appends the given middleware reference functions to the controller's
// list of middleware handlers. The middleware handlers are run in the order
// they are added to the controller. The middleware handlers are run before the
// controller's handlers. The controller's middleware handlers are run after the
// module's middleware handlers. The module middleware handlers are run after
// the module's parent middleware handlers. The module middleware handlers are
// run before the module's controllers. The controller middleware handlers are
// run before the controller's handlers. The middleware reference functions are
// called with the controller as the first argument and the ctx as the second
// argument. The middleware reference functions should return a middleware
// handler that takes a ctx as an argument and returns an error. If the
// middleware handler returns an error, the request is rejected with the
// error.
func (c *DynamicController) UseRef(middlewareRefs ...MiddlewareRef) *DynamicController {
	for _, v := range middlewareRefs {
		mid := c.ParseMiddlewareRef(v)
		c.middlewares = append(c.middlewares, mid)
	}
	return c
}

// Use appends the given middleware functions to the module's list of
// middleware handlers. The middleware handlers are run in the order they are
// added to the module. The middleware handlers are run before the module's
// controllers. The module middleware handlers are run after the module's
// parent middleware handlers. The module middleware handlers are run before
// the module's controllers. The module middleware handlers are run before the
// module's routers. The module middleware handlers are run before the module's
// handlers.
func (module *DynamicModule) Use(middlewares ...Middleware) *DynamicModule {
	module.Middlewares = append(module.Middlewares, middlewares...)
	for _, middleware := range middlewares {
		for _, router := range module.Routers {
			router.Middlewares = append(router.Middlewares, middleware)
		}
	}

	return module
}

func (module *DynamicModule) ParseMiddlewareRef(ref MiddlewareRef) Middleware {
	return func(ctx Ctx) error {
		return ref(module, ctx)
	}
}

func (module *DynamicModule) UseRef(middlewareRefs ...MiddlewareRef) *DynamicModule {
	for _, v := range middlewareRefs {
		mid := module.ParseMiddlewareRef(v)
		module.Middlewares = append(module.Middlewares, mid)
		for _, router := range module.Routers {
			router.Middlewares = append(router.Middlewares, mid)
		}
	}
	return module
}
