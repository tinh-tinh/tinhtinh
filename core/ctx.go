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

// Req returns the original http.Request from the client.
func (ctx *Ctx) Req() *http.Request {
	return ctx.r
}

// Res returns the original http.ResponseWriter from the server.
func (ctx *Ctx) Res() http.ResponseWriter {
	return ctx.w
}

// Headers returns the value of the given HTTP header from the client's request.
// If the header is not present, it returns an empty string.
func (ctx *Ctx) Headers(key string) string {
	return ctx.r.Header.Get(key)
}

// Cookies returns the named cookie provided in the request or
// nil if no such cookie is present.
func (ctx *Ctx) Cookies(key string) *http.Cookie {
	cookie, err := ctx.r.Cookie(key)
	if err != nil {
		return nil
	}
	return cookie
}

// SetCookie adds a Set-Cookie header to the response.
//
// The provided key is the name of the cookie, and the provided value is the
// value of the cookie. The maxAge argument specifies the maximum age of the
// cookie in seconds.
//
// The cookie is marked as HttpOnly, Secure, and SameSite=Lax.
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

// BodyParser is a helper to parse the request body into a given interface
func (ctx *Ctx) BodyParser(payload interface{}) error {
	err := json.NewDecoder(ctx.r.Body).Decode(payload)
	if err != nil {
		return err
	}
	return nil
}

// Body returns the request body as a given interface.
func (ctx *Ctx) Body() interface{} {
	return ctx.Get(InBody)
}

// Params returns the route parameters as a given interface.
func (ctx *Ctx) Params() interface{} {
	return ctx.Get(InPath)
}

// Queries returns the query parameters as a given interface.
func (ctx *Ctx) Queries() interface{} {
	return ctx.Get(InQuery)
}

// Param returns the value of the route parameter with the given key.
// If the parameter is not present, it returns an empty string.
func (ctx *Ctx) Param(key string) string {
	val := ctx.r.PathValue(key)
	return val
}

// Query returns the value of the query string parameter with the given key.
// If the parameter is not present, it returns an empty string.
func (ctx *Ctx) Query(key string) string {
	val := ctx.r.URL.Query().Get(key)
	return val
}

// QueryInt returns the value of the query string parameter with the given key as an integer.
// If the parameter is not present, it panics.
func (ctx *Ctx) QueryInt(key string) int {
	val := ctx.r.URL.Query().Get(key)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return intVal
}

// QueryBool returns the value of the query string parameter with the given key as a boolean.
// If the parameter is not present, it panics.
func (ctx *Ctx) QueryBool(key string) bool {
	val := ctx.r.URL.Query().Get(key)
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return boolVal
}

// Status sets the HTTP status code for the response.
//
// It returns itself to facilitate method chaining.
//
// Example:
//
//	ctx.Status(http.StatusOK).JSON(core.Map{"message": "Hello, World!"})
func (ctx *Ctx) Status(statusCode int) *Ctx {
	ctx.w.WriteHeader(statusCode)
	return ctx
}

type Map map[string]interface{}

// JSON sends the given data as a JSON response.
//
// The Content-Type of the response is set to "application/json".
//
// If there is an error while encoding the data, it panics.
func (ctx *Ctx) JSON(data Map) {
	ctx.w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(ctx.w).Encode(data)
	if err != nil {
		panic(err)
	}
}

// Get returns the value associated with the given key from the request context.
//
// If no value is associated with the key, it returns nil.
func (ctx *Ctx) Get(key any) interface{} {
	val := ctx.r.Context().Value(key)
	return val
}

// Set sets the value associated with the given key in the request context.
//
// If the given key already has a value associated with it in the request context,
// it will be overwritten.
//
// The value is associated with the given key in the context of the request,
// and can be retrieved using the Get method.
func (ctx *Ctx) Set(key interface{}, val interface{}) {
	ctx.r = ctx.r.WithContext(context.WithValue(ctx.r.Context(), key, val))
}

// NewCtx creates a new Ctx from the given http.ResponseWriter and *http.Request.
//
// It returns a new Ctx with the given http.ResponseWriter and *http.Request set
// as its fields.
//
// The returned Ctx can be used to call any of the methods on Ctx, such as
// Req, Res, Headers, etc.
func NewCtx(w http.ResponseWriter, r *http.Request) Ctx {
	return Ctx{
		w: w,
		r: r,
	}
}

// ParseCtx wraps a function that takes a Ctx argument and returns an
// http.HandlerFunc that can be used to handle an HTTP request.
//
// The returned http.HandlerFunc calls the given function with a new Ctx
// constructed from the given http.ResponseWriter and *http.Request, and then
// returns without doing anything else.
//
// The returned http.HandlerFunc can be used as a handler for an HTTP request.
func ParseCtx(ctxFnc func(ctx Ctx)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewCtx(w, r)
		ctxFnc(ctx)
	})
}
