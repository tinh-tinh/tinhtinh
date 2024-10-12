package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/middleware/cookie"
	"github.com/tinh-tinh/tinhtinh/middleware/storage"
)

type Ctx struct {
	r        *http.Request
	w        http.ResponseWriter
	handler  http.Handler
	metadata []*Metadata
	app      *App
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

func (ctx *Ctx) SignedCookie(key string, val ...string) (string, error) {
	s, ok := ctx.Get(cookie.SIGNED_COOKIE).(*cookie.SecureCookie)
	if !ok {
		common.InternalServerException(ctx.Res(), "failed to get signed cookie")
		return "", errors.New("failed to get signed cookie")
	}
	if len(val) > 0 {
		encoded, err := s.Encrypt(val[0])
		if err != nil {
			fmt.Println(err)
			common.InternalServerException(ctx.Res(), "failed to encode signed cookie")
			return "", errors.New("failed to encode signed cookie")
		}
		cookie := &http.Cookie{
			Name:     key,
			Value:    encoded,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(ctx.w, cookie)
		return val[0], nil
	}

	cookie, err := ctx.Req().Cookie(key)
	if err != nil {
		common.InternalServerException(ctx.Res(), "failed to get signed cookie")
		return "", errors.New("failed to get signed cookie")
	}
	value, err := s.Decrypt(cookie.Value)
	if err != nil {
		common.InternalServerException(ctx.Res(), "failed to decode signed cookie")
		return "", errors.New("failed to decode signed cookie")
	}
	return value, nil
}

// BodyParser is a helper to parse the request body into a given interface
func (ctx *Ctx) BodyParser(payload interface{}) error {
	body, err := io.ReadAll(ctx.r.Body)
	if err != nil {
		return err
	}

	err = ctx.app.decoder(body, payload)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) QueryParse(payload interface{}) error {
	ct := reflect.ValueOf(payload).Elem()
	fmt.Println(ct)
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("query")
		if tagVal != "" {
			val := ctx.Req().URL.Query().Get(tagVal)
			if val == "" {
				continue
			}
			switch field.Type.Name() {
			case "string":
				ct.Field(i).SetString(val)
			case "int":
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				ct.Field(i).SetInt(int64(intVal))
			case "bool":
				boolVal, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				ct.Field(i).SetBool(boolVal)
			default:
				return fmt.Errorf("unsupported type: %s", field.Type.Name())
			}
		}
	}
	return nil
}

func (ctx *Ctx) ParamParse(payload interface{}) error {
	ct := reflect.ValueOf(payload).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("param")
		if tagVal != "" {
			val := ctx.Req().PathValue(tagVal)
			if val == "" {
				continue
			}
			switch field.Type.Name() {
			case "string":
				ct.Field(i).SetString(val)
			case "int":
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				ct.Field(i).SetInt(int64(intVal))
			case "bool":
				boolVal, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				ct.Field(i).SetBool(boolVal)
			default:
				return fmt.Errorf("unsupported type: %s", field.Type.Name())
			}
		}
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
	res, err := ctx.app.encoder(data)
	if err != nil {
		panic(err)
	}

	_, err = ctx.w.Write(res)
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

func (ctx *Ctx) Next() error {
	ctx.handler.ServeHTTP(ctx.w, ctx.r)
	return nil
}

func (ctx *Ctx) Session(key string, val ...interface{}) interface{} {
	if len(val) > 0 {
		cookie := ctx.app.session.Set(key, val[0])
		http.SetCookie(ctx.w, &cookie)
		return nil
	}
	cookie, err := ctx.Req().Cookie(key)
	if err != nil {
		return nil
	}
	data := ctx.app.session.Get(cookie.Value)
	return data
}

// NewCtx creates a new Ctx from the given http.ResponseWriter and *http.Request.
//
// It returns a new Ctx with the given http.ResponseWriter and *http.Request set
// as its fields.
//
// The returned Ctx can be used to call any of the methods on Ctx, such as
// Req, Res, Headers, etc.
func NewCtx(app *App) *Ctx {
	return &Ctx{
		app: app,
	}
}

func (ctx *Ctx) SetCtx(w http.ResponseWriter, r *http.Request) {
	ctx.w = w
	ctx.r = r
}

func (ctx *Ctx) SetHandler(h http.Handler) {
	ctx.handler = h
}

// ParseCtx wraps a function that takes a Ctx argument and returns an
// http.HandlerFunc that can be used to handle an HTTP request.
//
// The returned http.HandlerFunc calls the given function with a new Ctx
// constructed from the given http.ResponseWriter and *http.Request, and then
// returns without doing anything else.
//
// The returned http.HandlerFunc can be used as a handler for an HTTP request.
func ParseCtx(app *App, ctxFnc func(ctx Ctx), meta ...*Metadata) http.Handler {
	ctx := app.pool.Get().(*Ctx)
	defer app.pool.Put(ctx)

	ctx.SetMetadata(meta...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx.SetCtx(w, r)
		ctxFnc(*ctx)
	})
}

func (ctx *Ctx) UploadedFile() *storage.File {
	uploadedFile, ok := ctx.Get(FILE).(*storage.File)
	if uploadedFile == nil || !ok {
		return nil
	}
	return uploadedFile
}

func (ctx *Ctx) UploadedFiles() []*storage.File {
	uploadedFiles, ok := ctx.Get(FILES).([]*storage.File)
	if uploadedFiles == nil || !ok {
		return nil
	}
	return uploadedFiles
}

func (ctx *Ctx) UploadedFieldFile() map[string][]*storage.File {
	uploadedFieldFile, ok := ctx.Get(FIELD_FILES).(map[string][]*storage.File)
	if uploadedFieldFile == nil || !ok {
		return nil
	}
	return uploadedFieldFile
}
