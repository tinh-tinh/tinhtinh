package api

import "net/http"

type Middleware func(http.Handler) http.Handler

// Guard func
type Guard func(ctx Ctx) bool

func ParseGuard(guard Guard) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(NewCtx(w, r))
			if !isAccess {
				http.Error(w, "You can not access", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// Interceptor
type Interceptor func(ctx Ctx)

// Pipe
type Pipe func(ctx Ctx)
