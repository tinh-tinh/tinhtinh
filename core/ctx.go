package core

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type Ctx struct {
	r    *http.Request
	w    http.ResponseWriter
	ctrl *DynamicController
}

func (ctx *Ctx) Req() *http.Request {
	return ctx.r
}

func (ctx *Ctx) Res() http.ResponseWriter {
	return ctx.w
}

func (ctx *Ctx) Headers(key string) string {
	return ctx.r.Header.Get(key)
}

func (ctx *Ctx) Cookies(key string) *http.Cookie {
	cookie, err := ctx.r.Cookie(key)
	if err != nil {
		return nil
	}
	return cookie
}

func (ctx *Ctx) SetCookie(key string, value string, maxAge int) {
	http.SetCookie(ctx.w, &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (ctx *Ctx) BodyParser(payload interface{}) error {
	err := json.NewDecoder(ctx.r.Body).Decode(payload)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) Body() interface{} {
	return ctx.Get(InBody)
}

func (ctx *Ctx) Params() interface{} {
	return ctx.Get(InPath)
}

func (ctx *Ctx) Queries() interface{} {
	return ctx.Get(InQuery)
}

func (ctx *Ctx) Param(key string) string {
	val := ctx.r.PathValue(key)
	return val
}

func (ctx *Ctx) Query(key string) string {
	val := ctx.r.URL.Query().Get(key)
	return val
}

func (ctx *Ctx) QueryInt(key string) int {
	val := ctx.r.URL.Query().Get(key)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return intVal
}

func (ctx *Ctx) QueryBool(key string) bool {
	val := ctx.r.URL.Query().Get(key)
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return boolVal
}
func (ctx *Ctx) Status(statusCode int) *Ctx {
	ctx.w.WriteHeader(statusCode)
	return ctx
}

type Map map[string]interface{}

func (ctx *Ctx) JSON(data Map) {
	ctx.w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(ctx.w).Encode(data)
	if err != nil {
		panic(err)
	}
}

func (ctx *Ctx) Get(key any) interface{} {
	val := ctx.r.Context().Value(key)
	return val
}

func (ctx *Ctx) Set(key interface{}, val interface{}) {
	ctx.r = ctx.r.WithContext(context.WithValue(ctx.r.Context(), key, val))
}

type InjectFnc func(req *http.Request) interface{}

func (ctx *Ctx) Inject(injFnc InjectFnc) interface{} {
	return injFnc(ctx.r)
}

func NewCtx(w http.ResponseWriter, r *http.Request) Ctx {
	return Ctx{
		w: w,
		r: r,
	}
}

func (ctrl *DynamicController) ParseCtx(ctxFnc func(ctx Ctx)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewCtx(w, r)
		ctxFnc(ctx)
	})
}
