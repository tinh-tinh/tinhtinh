package core

import (
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

// Guard is a function that checks access permission for a controller
type Guard func(ctx Ctx) bool

// ParseGuard wraps a Guard function into a Middleware that checks access permission
// for the given DynamicController. If the guard function returns false, it responds
// with a forbidden error, otherwise it calls the next middleware in the chain.
func (ctrl *DynamicController) ParseGuard(guard Guard) Middleware {
	return func(ctx Ctx) error {
		isAccess := guard(ctx)
		if !isAccess {
			return common.ForbiddenException(ctx.Res(), "you can not access")
		}
		return ctx.Next()
	}
}

// Guard registers the given Guard functions with the controller. The Guard functions
// are called in order, and if any of them return false, the request is rejected with a
// forbidden error. Otherwise, the request is allowed. Guard functions are called
// before the controller's middleware handlers. Guard functions are called after the
// module's middleware handlers. The module middleware handlers are run after the
// module's parent middleware handlers. The module middleware handlers are run before
// the module's controllers. The controller's Guard functions are run before the
// controller's handlers. If any of the Guard functions return an error, the request
// is rejected with the error.
func (c *DynamicController) Guard(guards ...Guard) Controller {
	for _, v := range guards {
		mid := c.ParseGuard(v)
		c.middlewares = append(c.middlewares, mid)
	}
	return c
}

// ParseGuard wraps an AppGuard function into a Middleware that checks access permission
// for the given DynamicModule. If the guard function returns false, it responds
// with a forbidden error, otherwise it calls the next middleware in the chain.
func (module *DynamicModule) ParseGuard(guard Guard) Middleware {
	return func(ctx Ctx) error {
		isAccess := guard(ctx)
		if !isAccess {
			return common.ForbiddenException(ctx.Res(), "you can not access")
		}
		return ctx.Next()
	}
}

func (module *DynamicModule) Guard(guards ...Guard) Module {
	for _, v := range guards {
		mid := module.ParseGuard(v)
		module.Middlewares = append(module.Middlewares, mid)
		for _, router := range module.Routers {
			router.Middlewares = append(router.Middlewares, mid)
		}
	}

	return module
}
