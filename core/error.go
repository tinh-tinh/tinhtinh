package core

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/tinh-tinh/tinhtinh/common"
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

type NotFoundHandler func(ctx Ctx) error

func NotFoundHandlerDefault(ctx Ctx) error {
	common.BadRequestException(ctx.Res(), fmt.Sprintf("not found: %s", ctx.Req().URL.Path))
	return nil
}

func checkRouter(app *App, r *http.Request) bool {
	idxRoute := slices.IndexFunc(app.Module.Routers, func(e *Router) bool {
		route := ParseRoute(r.Method + " " + r.URL.Path)
		if app.version != nil && app.version.Type == URIVersion && e.Version != "" {
			route.SetPrefix("v" + e.Version)
		}
		route.SetPrefix(e.Name)
		if app.Prefix != "" {
			route.SetPrefix(app.Prefix)
		}
		return route.GetPath() == r.URL.Path && r.Method == e.Method
	})
	return idxRoute != -1
}
