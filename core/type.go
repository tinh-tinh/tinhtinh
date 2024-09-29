package core

import (
	"net/http"
)

type Interceptor func(ctx Ctx) http.Handler

type Encode func(v interface{}) ([]byte, error)

type Decode func(data []byte, v interface{}) error
