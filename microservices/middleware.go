package microservices

import "fmt"

type Middleware func(ctx Ctx) error

type middlewareRaw func(m Factory) Factory

func (h *Handler) Use(middlewares ...Middleware) *Handler {
	h.middlewares = append(h.middlewares, middlewares...)
	return h
}

func (h *Handler) Registry() *Handler {
	h.globalMiddlewares = append(h.globalMiddlewares, h.middlewares...)
	h.middlewares = nil
	return h
}

func ParseCtxMiddleware(ctxMid Middleware) middlewareRaw {
	return func(f Factory) Factory {
		return FactoryFunc(func(ctx Ctx) error {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
					ctx.ErrorHandler(err)
				}
			}()
			err = ctxMid(ctx)
			if err != nil {
				ctx.ErrorHandler(err)
				return err
			}
			return f.Handle(ctx)
		})
	}
}
