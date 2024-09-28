package core

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Interceptor func(ctx Ctx) http.Handler

type Encode func(v interface{}) ([]byte, error)

type Decode func(data []byte, v interface{}) error
