package microservices

import (
	"net"

	"github.com/tinh-tinh/tinhtinh/core"
)

type HandlerFnc func(data interface{})

type Handler struct {
	core.DynamicController
	Conn net.Conn
}

func NewHandler(ctrl *core.DynamicController, conn net.Conn) *Handler {
	return &Handler{DynamicController: *ctrl, Conn: conn}
}

func (h *Handler) EventPattern(event string, cb HandlerFnc) {

}
