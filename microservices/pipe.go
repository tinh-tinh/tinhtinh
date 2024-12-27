package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

type CtxKey string

const PIPE CtxKey = "pipe"

func PipeMiddleware(value interface{}) Middleware {
	return func(ctx Ctx) error {
		schema := ctx.Payload(value)
		err := validator.Scanner(schema)
		if err != nil {
			panic(exception.ThrowRpc(err.Error()))
		}
		ctx.Set(PIPE, schema)
		return ctx.Next()
	}
}

func (h *Handler) Pipe(value interface{}) *Handler {
	h.middlewares = append(h.middlewares, PipeMiddleware(value))
	return h
}
