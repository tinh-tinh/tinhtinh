package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Guard func(ref core.RefProvider, ctx Ctx) bool

func (h *Handler) ParseGuard(guard Guard) Middleware {
	return func(ctx Ctx) error {
		isAccess := guard(h, ctx)
		if !isAccess {
			panic(exception.ThrowRpc("the service reject connect"))
		}
		return ctx.Next()
	}
}

func (h *Handler) Guard(guards ...Guard) *Handler {
	for _, g := range guards {
		h.middlewares = append(h.middlewares, h.ParseGuard(g))
	}
	return h
}
