package core

import (
	"net/http"
)

// Guard func
type Guard func(module *DynamicModule, ctx Ctx) bool

func (module *DynamicModule) ParseGuard(guard Guard) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(module, NewCtx(w, r))
			if !isAccess {
				ForbiddenException(w, "you can not access")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

type GuardWithCtrl func(ctrl *DynamicController, ctx Ctx) bool

func (ctrl *DynamicController) ParseGuardCtrl(guard GuardWithCtrl) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(ctrl, NewCtx(w, r))
			if !isAccess {
				ForbiddenException(w, "you can not access")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
