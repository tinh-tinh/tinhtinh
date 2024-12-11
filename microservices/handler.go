package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Handler struct {
	core.DynamicProvider
	module core.Module
}

func NewHandler(module core.Module, opt core.ProviderOptions) *Handler {
	provider := &Handler{}
	provider.module = module

	return provider
}

func (h *Handler) OnResponse(name string, fnc core.Factory) {
	core.InitProviders(h.module, core.ProviderOptions{Name: core.Provide(name), Factory: fnc, Scope: h.Scope})
}

func (h *Handler) OnEvent(name string, fnc core.Factory) {
	core.InitProviders(h.module, core.ProviderOptions{Name: core.Provide(name), Factory: fnc, Scope: h.Scope})
}
