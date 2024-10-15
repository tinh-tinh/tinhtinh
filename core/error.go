package core

import (
	"net/http"
	"time"
)

type ErrorHandler func(err error, ctx Ctx) error

// ErrorHandlerDefault is the default error handler of the app.
//
// It will return a JSON response with a status code of 500, containing the
// timestamp and the path of the request.
//
// If the error is nil, it will return nil without doing anything.
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