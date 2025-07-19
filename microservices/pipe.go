package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
)

type CtxKey string

const PIPE CtxKey = "pipe"

func PipeMiddleware(dto PipeDto) Middleware {
	return func(ctx Ctx) error {
		payload := dto.GetValue()
		err := ctx.PayloadParser(payload)
		if err != nil {
			return exception.ThrowRpc(err.Error())
		}

		err = ctx.Scan(payload)
		if err != nil {
			return exception.ThrowRpc(err.Error())
		}
		ctx.Set(PIPE, payload)
		return ctx.Next()
	}
}

func (h *Handler) Pipe(value PipeDto) *Handler {
	h.middlewares = append(h.middlewares, PipeMiddleware(value))
	return h
}

type PipeDto interface {
	GetValue() interface{}
}

type PayloadParser[P any] struct {
}

func (p PayloadParser[P]) GetValue() interface{} {
	var payload P
	return &payload
}
