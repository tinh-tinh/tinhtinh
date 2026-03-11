package microservices

import "fmt"

type EventFactory interface {
	Handle(ctx Ctx) error
}

type SubscribeHandler struct {
	Name        string
	Factory     EventFactory
	Middlewares []Middleware
}

type EventFactoryFunc func(ctx Ctx) error

func (f EventFactoryFunc) Handle(ctx Ctx) error {
	return f(ctx)
}

func ParseEventFactory(factory EventFactory) EventFactory {
	return EventFactoryFunc(func(ctx Ctx) error {
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
				ctx.ErrorHandler(err)
			}
		}()

		err = factory.Handle(ctx)
		if err != nil {
			ctx.ErrorHandler(err)
			return err
		}
		return nil
	})
}

func (s *SubscribeHandler) Handle(svc Service, data Message) error {
	mergeHandler := ParseEventFactory(s.Factory)

	for i := len(s.Middlewares) - 1; i >= 0; i-- {
		v := s.Middlewares[i]
		mid := ParseCtxMiddleware(v)
		mergeHandler = mid(mergeHandler)
	}

	return mergeHandler.Handle(NewCtx(data, svc))
}
