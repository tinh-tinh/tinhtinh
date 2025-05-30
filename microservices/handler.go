package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func NewHandler(module core.Module, opt core.ProviderOptions) *Handler {
	provider := &Handler{}
	provider.module = module

	return provider
}

type Handler struct {
	core.DynamicProvider
	module            core.Module
	middlewares       []Middleware
	globalMiddlewares []Middleware
}

// OnResponse registers a provider with the given name and factory function to be
// called when the response is ready. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnResponse(name string, fnc FactoryFunc) {
	core.InitProviders(h.module, core.ProviderOptions{
		Name: core.Provide(name),
		Factory: func(param ...interface{}) interface{} {
			store, ok := param[0].(*Store)
			if !ok {
				return nil
			}
			if store.Subscribers[RPC] == nil {
				store.Subscribers[RPC] = []*SubscribeHandler{}
			}
			store.Subscribers[RPC] = append(store.Subscribers[RPC], &SubscribeHandler{
				Name:        name,
				Factory:     fnc,
				Middlewares: append(h.globalMiddlewares, h.middlewares...),
			})
			return store.Subscribers
		},
		Inject: []core.Provide{STORE},
		Scope:  h.Scope,
	})
	h.middlewares = nil
}

// OnEvent registers a provider with the given name and factory function to be
// called when an event is triggered. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnEvent(name string, fnc FactoryFunc) {
	core.InitProviders(h.module, core.ProviderOptions{
		Name: core.Provide(name),
		Factory: func(param ...interface{}) interface{} {
			store, ok := param[0].(*Store)
			if !ok {
				return nil
			}
			if store.Subscribers[PubSub] == nil {
				store.Subscribers[PubSub] = []*SubscribeHandler{}
			}
			store.Subscribers[PubSub] = append(store.Subscribers[PubSub], &SubscribeHandler{
				Name:        name,
				Factory:     fnc,
				Middlewares: append(h.globalMiddlewares, h.middlewares...),
			})
			return store.Subscribers
		},
		Inject: []core.Provide{STORE}, Scope: h.Scope,
	})
	h.middlewares = nil
}

func (h *Handler) Ref(name core.Provide, ctx ...core.Ctx) interface{} {
	return h.module.Ref(name, ctx...)
}
