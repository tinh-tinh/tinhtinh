package core

import (
	"net/http"
)

// Guard func
type Guard func(ctx Ctx) bool

func ParseGuard(guard Guard) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(NewCtx(w, r))
			if !isAccess {
				ForbiddenException(w, "you can not access")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

type GuardWithCtrl func(ctrl *DynamicController, ctx Ctx) bool

func ParseGuardCtrl(ctrl *DynamicController, guard GuardWithCtrl) Middleware {
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
