package core

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Interceptor func(ctx Ctx) http.Handler
