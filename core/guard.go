package core

import (
	"errors"

	"github.com/tinh-tinh/tinhtinh/common"
)

type Guard func(ctrl *DynamicController, ctx Ctx) bool

// ParseGuard takes a GuardWithCtrl function and returns a Middleware that checks
// if the user has access according to the GuardWithCtrl function. If the user does
// not have access, it returns a 403 status code and ends the request.
//
// The GuardWithCtrl function is passed the current controller and a new Ctx as arguments.
func (ctrl *DynamicController) ParseGuard(guard Guard) Middleware {
	return func(ctx Ctx) error {
		isAccess := guard(ctrl, ctx)
		if !isAccess {
			common.ForbiddenException(ctx.Res(), "you can not access")
			return errors.New("you can not access")
		}
		return ctx.Next()
	}
}
