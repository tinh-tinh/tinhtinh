package core

import (
	"github.com/tinh-tinh/tinhtinh/common"
)

// GuardWithCtrl is a function that takes a controller and a Ctx as arguments
// and returns a bool.
type Guard func(ctrl *DynamicController, ctx *Ctx) bool

// ParseGuard takes a GuardWithCtrl function and returns a Middleware that checks
// if the user has access according to the GuardWithCtrl function. If the user does
// not have access, it returns a 403 status code and ends the request.
//
// The GuardWithCtrl function is passed the current controller and a new Ctx as arguments.
func (ctrl *DynamicController) ParseGuard(guard Guard) Middleware {
	return func(ctx Ctx) error {
		isAccess := guard(ctrl, &ctx)
		if !isAccess {
			return common.ForbiddenException(ctx.Res(), "you can not access")
		}
		return ctx.Next()
	}
}

// Guard registers the given guard functions with the controller. The guard
// functions are run in the order they are added to the controller. The guard
// functions are run before the controller's middleware handlers. The guard
// functions are run after the module's middleware handlers. The guard functions
// are run before the controller's handlers. If any of the guard functions
// return false, the request will be rejected with a 403 status code.
func (c *DynamicController) Guard(guards ...Guard) *DynamicController {
	for _, v := range guards {
		mid := c.ParseGuard(v)
		c.middlewares = append(c.middlewares, mid)
	}
	return c
}

// AppGuard is a function that takes a module and a Ctx as arguments and returns
// true if the user has access to the module, false otherwise.
type AppGuard func(module *DynamicModule, ctx Ctx) bool

// ParseGuard takes an AppGuard function and returns a Middleware that checks
// if the user has access according to the AppGuard function. If the user does
// not have access, it returns a 403 status code and ends the request.
//
// The AppGuard function is passed the current module and a new Ctx as arguments.
func (module *DynamicModule) ParseGuard(guard AppGuard) Middleware {
	return func(ctx Ctx) error {
		isAccess := guard(module, ctx)
		if !isAccess {
			return common.ForbiddenException(ctx.Res(), "you can not access")
		}
		return ctx.Next()
	}
}
