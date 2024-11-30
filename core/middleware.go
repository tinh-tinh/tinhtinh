package core

import (
	"fmt"
	"net/http"
)

type Middleware func(ctx Ctx) error

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
