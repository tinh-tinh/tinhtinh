package tcp

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

type RpcGateway struct {
	handlers []*microservices.RpcHandler
	service  microservices.Service
}

func (g *RpcGateway) Call(request *[]byte, reply *[]byte) error {
	msg := microservices.DecodeMessage(g.service, *request)
	handler, found := common.Find(g.handlers, func(h *microservices.RpcHandler) bool {
		return h.Name == msg.Event
	})
	if !found {
		return fmt.Errorf("handler not found")
	}

	handlerFnc := handler.Factory
	for i := len(handler.Middlewares) - 1; i >= 0; i-- {
		mid := handler.Middlewares[i]
		next := handlerFnc
		handlerFnc = func(ctx microservices.Ctx) ([]byte, error) {
			err := mid(ctx)
			if err != nil {
				return nil, err
			}
			return next(ctx)
		}
	}

	safeHandlerFnc := func(ctx microservices.Ctx) (res []byte, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
				ctx.ErrorHandler(err)
			}
		}()
		return handlerFnc(ctx)
	}

	ctx := microservices.NewCtx(msg, g.service)
	res, err := safeHandlerFnc(ctx)
	if err != nil {
		return err
	}

	*reply = res
	return nil
}
