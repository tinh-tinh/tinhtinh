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

const STORE core.Provide = "STORE"

type Store struct {
	Subscribers map[string][]SubscribeHandler
}

func Register() core.Modules {
	return func(module core.Module) core.Module {
		handlerModule := module.New(core.NewModuleOptions{})

		handlerModule.NewProvider(core.ProviderOptions{
			Name:  STORE,
			Value: &Store{Subscribers: make(map[string][]SubscribeHandler)},
		})

		handlerModule.Export(STORE)
		return handlerModule
	}
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
			if store.Subscribers[string(RPC)] == nil {
				store.Subscribers[string(RPC)] = []SubscribeHandler{}
			}
			store.Subscribers[string(RPC)] = append(store.Subscribers[string(RPC)], SubscribeHandler{
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
			if store.Subscribers[string(PubSub)] == nil {
				store.Subscribers[string(PubSub)] = []SubscribeHandler{}
			}
			store.Subscribers[string(PubSub)] = append(store.Subscribers[string(PubSub)], SubscribeHandler{
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
