package core

type Wrapper[V any] struct {
	handler func(data interface{}, ctx Ctx) V
}

type WrappedCtx[V any] struct {
	Ctx  Ctx
	Data V
}

func CreateWrapper[V any](handler func(data interface{}, ctx Ctx) V) Wrapper[V] {
	return Wrapper[V]{
		handler: handler,
	}
}

func (w *Wrapper[V]) Handler(data interface{}, fnc func(w WrappedCtx[V]) error) Handler {
	return func(ctx Ctx) error {
		wrappedCtx := WrappedCtx[V]{
			Ctx:  ctx,
			Data: w.handler(data, ctx),
		}

		return fnc(wrappedCtx)
	}
}
