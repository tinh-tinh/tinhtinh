package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Handler struct {
	core.DynamicProvider
	module core.Module
	schema interface{}
}

// NewHandler creates a new Handler with the given module and options.
// It returns the created Handler.
func NewHandler(module core.Module, opt core.ProviderOptions) *Handler {
	provider := &Handler{}
	provider.module = module

	return provider
}

func (h *Handler) Schema(schema interface{}) *Handler {
	h.schema = schema
	return h
}

// OnResponse registers a provider with the given name and factory function to be
// called when the response is ready. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnResponse(name string, fnc Factory) {
	core.InitProviders(h.module, core.ProviderOptions{Name: core.Provide(name), Factory: ConvertFactory(fnc), Scope: h.Scope, Type: core.EVENT})
}

// OnEvent registers a provider with the given name and factory function to be
// called when an event is triggered. The provider will be registered with the
// same scope as the handler.
func (h *Handler) OnEvent(name string, fnc Factory) {
	core.InitProviders(h.module, core.ProviderOptions{Name: core.Provide(name), Factory: ConvertFactory(fnc), Scope: h.Scope, Type: core.EVENT})
}
