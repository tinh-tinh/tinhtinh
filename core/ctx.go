package core

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type Ctx struct {
	r *http.Request
	w http.ResponseWriter
}

func (ctx *Ctx) BodyParser(payload interface{}) error {
	err := json.NewDecoder(ctx.r.Body).Decode(payload)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) Params(key string) string {
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

func (ctx *Ctx) Input() interface{} {
	val := ctx.Get(Input)
	return val
}

func (ctx *Ctx) Get(key any) interface{} {
	val := ctx.r.Context().Value(key)
	return val
}

func (ctx *Ctx) Set(key interface{}, val interface{}) {
	ctx.r = ctx.r.WithContext(context.WithValue(ctx.r.Context(), key, val))
}

func NewCtx(w http.ResponseWriter, r *http.Request) Ctx {
	return Ctx{
		w: w,
		r: r,
	}
}

func ParseCtx(ctxFnc func(ctx Ctx)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewCtx(w, r)
		ctxFnc(ctx)
	})
}
