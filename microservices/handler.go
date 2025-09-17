package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func NewHandler(module core.Module, transports ...string) *Handler {
	provider := &Handler{}
	provider.module = module

	for _, transport := range transports {
		provider.transports = append(provider.transports, ToTransport(transport))
	}

	return provider
}

type Handler struct {
	core.DynamicProvider
	module            core.Module
	middlewares       []Middleware
	globalMiddlewares []Middleware
	transports        []core.Provide
}

// OnResponse registers a provider with the given name and factory function to be
// called when the response is ready. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnResponse(name string, fnc FactoryFunc) {
	refNames := []core.Provide{STORE}
	refNames = append(refNames, h.transports...)

	for _, refName := range refNames {
		store, ok := h.module.Ref(refName).(*Store)
		if !ok {
			return
		}
		if store.Subscribers[RPC] == nil {
			store.Subscribers[RPC] = []*SubscribeHandler{}
		}
		store.Subscribers[RPC] = append(store.Subscribers[RPC], &SubscribeHandler{
			Name:        name,
			Factory:     fnc,
			Middlewares: append(h.globalMiddlewares, h.middlewares...),
		})
	}

	h.middlewares = nil
}

// OnEvent registers a provider with the given name and factory function to be
// called when an event is triggered. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnEvent(name string, fnc FactoryFunc) {
	refNames := []core.Provide{STORE}
	refNames = append(refNames, h.transports...)

	for _, refName := range refNames {
		store, ok := h.module.Ref(refName).(*Store)
		if !ok {
			return
		}
		if store.Subscribers[PubSub] == nil {
			store.Subscribers[PubSub] = []*SubscribeHandler{}
		}
		store.Subscribers[PubSub] = append(store.Subscribers[PubSub], &SubscribeHandler{
			Name:        name,
			Factory:     fnc,
			Middlewares: append(h.globalMiddlewares, h.middlewares...),
		})
	}
	h.middlewares = nil
}

func (h *Handler) Ref(name core.Provide, ctx ...core.Ctx) interface{} {
	return h.module.Ref(name, ctx...)
}
