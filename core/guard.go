package core

import (
	"net/http"

	"github.com/tinh-tinh/tinhtinh/common"
)

// Guard func
type Guard func(module *DynamicModule, ctx Ctx) bool

// ParseGuard takes a Guard function and returns a Middleware that checks
// if the user has access according to the Guard function. If the user does
// not have access, it returns a 403 status code and ends the request.
func (module *DynamicModule) ParseGuard(guard Guard) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(module, NewCtx(w, r))
			if !isAccess {
				common.ForbiddenException(w, "you can not access")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

type GuardWithCtrl func(ctrl *DynamicController, ctx Ctx) bool

// ParseGuardCtrl takes a GuardWithCtrl function and returns a Middleware that checks
// if the user has access according to the GuardWithCtrl function. If the user does
// not have access, it returns a 403 status code and ends the request.
//
// The GuardWithCtrl function is passed the current controller and a new Ctx as arguments.
func (ctrl *DynamicController) ParseGuardCtrl(guard GuardWithCtrl) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(ctrl, NewCtx(w, r))
			if !isAccess {
				common.ForbiddenException(w, "you can not access")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
