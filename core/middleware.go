package core

import "net/http"

type Middleware func(ctx Ctx) error

type middlewareRaw func(http.Handler) http.Handler

func ParseCtxMiddleware(app *App, ctxMid Middleware) middlewareRaw {
	ctx := app.pool.Get().(*Ctx)
	defer app.pool.Put(ctx)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.SetCtx(w, r)
			ctx.SetHandler(h)
			err := ctxMid(*ctx)
			if err != nil {
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

func (module *DynamicModule) Use(middlewares ...Middleware) *DynamicModule {
	module.Middlewares = append(module.Middlewares, middlewares...)
	for _, middleware := range middlewares {
		for _, router := range module.Routers {
			router.Middlewares = append(router.Middlewares, middleware)
		}
	}

	return module
}
