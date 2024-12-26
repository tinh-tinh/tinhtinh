package microservices

import "fmt"

type Factory interface {
	Handle(ctx Ctx) error
}

type SubscribeHandler struct {
	Name        string
	Factory     Factory
	Middlewares []Middleware
}

type FactoryFunc func(ctx Ctx) error

func (f FactoryFunc) Handle(ctx Ctx) error {
	return f(ctx)
}

func ParseFactory(factory Factory) Factory {
	return FactoryFunc(func(ctx Ctx) error {
		fmt.Println("Ctx is ", ctx)
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
				ctx.ErrorHandler(err)
			}
		}()

		err = factory.Handle(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *SubscribeHandler) Handle(svc Service, data interface{}) {
	var mergeHandler Factory
	mergeHandler = ParseFactory(s.Factory)

	for i := len(s.Middlewares) - 1; i >= 0; i-- {
		v := s.Middlewares[i]
		mid := ParseCtxMiddleware(v)
		mergeHandler = mid(mergeHandler)
	}

	mergeHandler.Handle(NewCtx(data, svc))
}
