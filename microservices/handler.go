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

// OnEvent registers a provider with the given name and factory function to be
// called when an event is triggered. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnEvent(name string, fnc FactoryFunc) {
	var refNames []core.Provide
	if len(h.transports) > 0 {
		refNames = append(refNames, h.transports...)
	} else {
		refNames = []core.Provide{STORE}
	}

	for _, refName := range refNames {
		store, ok := h.module.Ref(refName).(*Store)
		if !ok {
			return
		}
		store.Subscribers = append(store.Subscribers, &SubscribeHandler{
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
