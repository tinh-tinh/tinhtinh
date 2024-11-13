package core

type Wrapper[I any, V any] struct {
	handler func(data I, ctx Ctx) V
}

type WrappedCtx[V any] struct {
	Ctx
	Data V
}

func CreateWrapper[I any, V any](handler func(data I, ctx Ctx) V) Wrapper[I, V] {
	return Wrapper[I, V]{
		handler: handler,
	}
}

func (w *Wrapper[I, V]) Handler(data I, fnc func(w WrappedCtx[V]) error) Handler {
	return func(ctx Ctx) error {
		wrappedCtx := WrappedCtx[V]{
			Ctx:  ctx,
			Data: w.handler(data, ctx),
		}

		return fnc(wrappedCtx)
	}
}
