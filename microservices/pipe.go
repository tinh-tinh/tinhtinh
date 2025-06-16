package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

type CtxKey string

const PIPE CtxKey = "pipe"

func PipeMiddleware(dto core.PipeDto) Middleware {
	return func(ctx Ctx) error {
		payload := dto.GetValue()
		err := ctx.PayloadParser(payload)
		if err != nil {
			return exception.ThrowRpc(err.Error())
		}

		err = validator.Scanner(payload)
		if err != nil {
			return exception.ThrowRpc(err.Error())
		}
		ctx.Set(PIPE, payload)
		return ctx.Next()
	}
}

func (h *Handler) Pipe(value core.PipeDto) *Handler {
	h.middlewares = append(h.middlewares, PipeMiddleware(value))
	return h
}

func Payload[P any](dto P) core.PipeDto {
	return &core.Pipe[P]{
		In: core.InBody,
	}
}
