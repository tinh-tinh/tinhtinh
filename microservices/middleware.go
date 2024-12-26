package microservices

import "fmt"

type Middleware func(ctx Ctx) error

type middlewareRaw func(m Factory) Factory

func (h *Handler) Use(middlewares ...Middleware) *Handler {
	h.middlewares = append(h.middlewares, middlewares...)
	return h
}

func (h *Handler) Registry() {
	h.globalMiddlewares = append(h.globalMiddlewares, h.middlewares...)
	h.middlewares = nil
}

func ParseCtxMiddleware(ctxMid Middleware) middlewareRaw {
	return func(f Factory) Factory {
		return FactoryFunc(func(ctx Ctx) error {
			ctx.SetFactory(f)
			fmt.Println("Ctx is ", ctx)
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
					ctx.ErrorHandler(err)
				}
			}()
			err = ctxMid(ctx)
			if err != nil {
				return err
			}
			return nil
		})
	}
}
