package core

import (
	"fmt"
	"net/http"
)

type (
	middlewareRaw func(http.Handler) http.Handler
	Middleware    func(ctx Ctx) error
)

type RefProvider interface {
	Ref(name Provide, ctx ...Ctx) interface{}
}

// ParseCtxMiddleware wraps a Middleware function and returns a middlewareRaw
// that can be used by http server. It provides a Ctx instance to the wrapped
// middleware function and automatically sets the handler of the Ctx instance.
func ParseCtxMiddleware(app *App, ctxMid Middleware, router *Router) middlewareRaw {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := app.pool.Get().(*DefaultCtx)
			defer app.pool.Put(ctx)
			ctx.SetMetadata(router.Metadata...)
			ctx.SetCtx(w, r)
			ctx.SetHandler(h)
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
					app.errorHandler(err, ctx)
				}
			}()
			err = ctxMid(ctx)
			if err != nil {
				app.errorHandler(err, ctx)
				return
			}
		})
	}
}

func (c *DynamicController) Use(middleware ...Middleware) Controller {
	c.middlewares = append(c.middlewares, middleware...)
	return c
}

func (module *DynamicModule) Use(middlewareRefs ...Middleware) Module {
	for _, mid := range middlewareRefs {
		module.Middlewares = append(module.Middlewares, mid)
		for _, router := range module.Routers {
			router.Middlewares = append(router.Middlewares, mid)
		}
	}
	return module
}
