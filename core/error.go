package core

import (
	"net/http"
	"time"
)

type ErrorHandler func(err error, ctx Ctx) error

func ErrorHandlerDefault(err error, ctx Ctx) error {
	if err != nil {
		res := Map{
			"statusCode": http.StatusInternalServerError,
			"timestamp":  time.Now().Format(time.RFC3339),
			"path":       ctx.Req().URL.Path,
		}
		return ctx.Status(http.StatusInternalServerError).JSON(res)
	}
	return nil
}
