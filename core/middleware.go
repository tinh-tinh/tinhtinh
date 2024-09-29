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
