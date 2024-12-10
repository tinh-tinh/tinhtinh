package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/tinh-tinh/tinhtinh/v2/middleware/cookie"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

type Ctx interface {
	Req() *http.Request
	Res() http.ResponseWriter
	Headers(key string) string
	Cookies(key string) *http.Cookie
	SetCookie(key string, value string, maxAge int)
	SignedCookie(key string, val ...string) (string, error)
	BodyParser(payload interface{}) error
	QueryParser(payload interface{}) error
	ParamParser(payload interface{}) error
	Body() interface{}
	Params() interface{}
	Queries() interface{}
	Param(key string) string
	ParamInt(key string) int
	ParamFloat(key string) float64
	ParamBool(key string) bool
	Query(key string) string
	QueryInt(key string) int
	QueryFloat(key string) float64
	QueryBool(key string) bool
	SetCallHandler(call CallHandler)
	JSON(data Map) error
	Get(key interface{}) interface{}
	Set(key interface{}, val interface{})
	Next() error
	Session(key string, val ...interface{}) interface{}
	SetCtx(w http.ResponseWriter, r *http.Request)
	SetHandler(h http.Handler)
	UploadedFile() *storage.File
	UploadedFiles() []*storage.File
	UploadedFieldFile() map[string][]*storage.File
	Redirect(uri string) error
	Ref(name Provide) interface{}
	GetMetadata(key string) interface{}
	ExportCSV(name string, body [][]string) error
	Status(statusCode int) Ctx
}

type DefaultCtx struct {
	r           *http.Request
	w           http.ResponseWriter
	handler     http.Handler
	metadata    []*Metadata
	callHandler CallHandler
	app         *App
}

// Req returns the original http.Request from the client.
func (ctx *DefaultCtx) Req() *http.Request {
	return ctx.r
}

// Res returns the original http.ResponseWriter from the server.
func (ctx *DefaultCtx) Res() http.ResponseWriter {
	return ctx.w
}

// Headers returns the value of the given HTTP header from the client's request.
// If the header is not present, it returns an empty string.
func (ctx *DefaultCtx) Headers(key string) string {
	return ctx.r.Header.Get(key)
}

// Cookies returns the named cookie provided in the request or
// nil if no such cookie is present.
func (ctx *DefaultCtx) Cookies(key string) *http.Cookie {
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
func (ctx *DefaultCtx) SetCookie(key string, value string, maxAge int) {
	http.SetCookie(ctx.w, &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// SignedCookie sets a signed cookie in the response, or gets a signed cookie
// from the request. If the signed cookie is set, it is encrypted using the
// app's secure cookie secret. If the signed cookie is retrieved, it is
// decrypted using the app's secure cookie secret. If the signed cookie is
// invalid, an error is returned.
//
// The first argument is the name of the cookie. The second argument is the
// value of the cookie to be set, or the empty string if the value is to be
// retrieved.
//
// The cookie is marked as HttpOnly, Secure, and SameSite=Lax.
func (ctx *DefaultCtx) SignedCookie(key string, val ...string) (string, error) {
	s, ok := ctx.Get(cookie.SIGNED_COOKIE).(*cookie.SecureCookie)
	if !ok {
		return "", errors.New("failed to get signed cookie")
	}
	if len(val) > 0 {
		encoded, err := s.Encrypt(val[0])
		if err != nil {
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
		return "", errors.New("failed to get signed cookie")
	}
	value, err := s.Decrypt(cookie.Value)
	if err != nil {
		return "", errors.New("failed to decode signed cookie")
	}
	return value, nil
}

// BodyParser is a helper to parse the request body into a given interface
func (ctx *DefaultCtx) BodyParser(payload interface{}) error {
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

// QueryParser takes a struct and populates its fields based on the query
// parameters in the request. It supports string, int, and bool types.
// If a field has a "query" tag, it will be populated with the query
// parameter of the same name. If the query parameter is not present,
// the field will be skipped.
//
// For example:
//
//	type MyStruct {
//		Name string `query:"name"`
//	}
//
//	req, _ := http.NewRequest("GET", "/?name=John", nil)
//	ctx := NewCtx(req, nil)
//	var ms MyStruct
//	err := ctx.QueryParser(&ms)
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(ms.Name) // John
func (ctx *DefaultCtx) QueryParser(payload interface{}) error {
	ct := reflect.ValueOf(payload).Elem()
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

// ParamParser takes a struct and populates its fields based on the path
// parameters in the request. It supports string, int, and bool types.
// If a field has a "param" tag, it will be populated with the path
// parameter of the same name. If the path parameter is not present,
// the field will be skipped.
func (ctx *DefaultCtx) ParamParser(payload interface{}) error {
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
func (ctx *DefaultCtx) Body() interface{} {
	return ctx.Get(InBody)
}

// Params returns the route parameters as a given interface.
func (ctx *DefaultCtx) Params() interface{} {
	return ctx.Get(InPath)
}

// Queries returns the query parameters as a given interface.
func (ctx *DefaultCtx) Queries() interface{} {
	return ctx.Get(InQuery)
}

// Param returns the value of the route parameter with the given key.
// If the parameter is not present, it returns an empty string.
func (ctx *DefaultCtx) Param(key string) string {
	val := ctx.r.PathValue(key)
	return val
}

// ParamInt returns the value of the route parameter with the given key as an integer.
// If the parameter is not present, it panics.
func (ctx *DefaultCtx) ParamInt(key string) int {
	val := ctx.r.PathValue(key)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return intVal
}

// ParamFloat returns the value of the route parameter with the given key as a float64.
// If the parameter is not present, it panics.
func (ctx *DefaultCtx) ParamFloat(key string) float64 {
	val := ctx.r.PathValue(key)
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return floatVal
}

// ParamBool returns the value of the route parameter with the given key as a boolean.
// If the parameter is not present, it panics.
func (ctx *DefaultCtx) ParamBool(key string) bool {
	val := ctx.r.PathValue(key)
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return boolVal
}

// Query returns the value of the query string parameter with the given key.
// If the parameter is not present, it returns an empty string.
func (ctx *DefaultCtx) Query(key string) string {
	val := ctx.r.URL.Query().Get(key)
	return val
}

// QueryInt returns the value of the query string parameter with the given key as an integer.
// If the parameter is not present, it panics.
func (ctx *DefaultCtx) QueryInt(key string) int {
	val := ctx.r.URL.Query().Get(key)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return intVal
}

// QueryFloat returns the value of the query string parameter with the given key as a float.
// If the parameter is not present, it panics.
func (ctx *DefaultCtx) QueryFloat(key string) float64 {
	val := ctx.r.URL.Query().Get(key)
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return floatVal
}

// QueryBool returns the value of the query string parameter with the given key as a boolean.
// If the parameter is not present, it panics.
func (ctx *DefaultCtx) QueryBool(key string) bool {
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
func (ctx *DefaultCtx) Status(statusCode int) Ctx {
	ctx.w.WriteHeader(statusCode)
	return ctx
}

type Map map[string]interface{}

func (ctx *DefaultCtx) SetCallHandler(call CallHandler) {
	ctx.callHandler = call
}

// JSON sends the given data as a JSON response.
//
// The Content-Type of the response is set to "application/json".
//
// If there is an error while encoding the data, it panics.
func (ctx *DefaultCtx) JSON(data Map) error {
	ctx.w.Header().Set("Content-Type", "application/json")
	if ctx.callHandler != nil {
		data = ctx.callHandler(data)
	}
	res, err := ctx.app.encoder(data)
	if err != nil {
		return err
	}

	_, err = ctx.w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

// Get returns the value associated with the given key from the request context.
//
// If no value is associated with the key, it returns nil.
func (ctx *DefaultCtx) Get(key any) interface{} {
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
func (ctx *DefaultCtx) Set(key interface{}, val interface{}) {
	ctx.r = ctx.r.WithContext(context.WithValue(ctx.r.Context(), key, val))
}

// Next calls the next handler in the chain, using the
// http.ResponseWriter and *http.Request from the current context.
//
// If there is no next handler, it does nothing.
//
// It returns nil.
func (ctx *DefaultCtx) Next() error {
	ctx.handler.ServeHTTP(ctx.w, ctx.r)
	return nil
}

// Session sets a session cookie or gets a session value.
//
// If a single argument is given, it sets a session cookie with the given key and value.
// The cookie is marked as HttpOnly, Secure, and SameSite=Lax.
//
// If no arguments are given, it gets the session value associated with the given key.
// If the cookie is not present, or if the value is not found, it returns nil.
func (ctx *DefaultCtx) Session(key string, val ...interface{}) interface{} {
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
func NewCtx(app *App) *DefaultCtx {
	return &DefaultCtx{
		app: app,
	}
}

// SetCtx sets the http.ResponseWriter and *http.Request fields of the Ctx
// to the given values.
//
// It returns nothing.
func (ctx *DefaultCtx) SetCtx(w http.ResponseWriter, r *http.Request) {
	ctx.w = w
	ctx.r = r
}

// SetHandler sets the http.Handler field of the Ctx to the given value.
//
// This is usually called by the framework to set the handler for the current request.
func (ctx *DefaultCtx) SetHandler(h http.Handler) {
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
func ParseCtx(app *App, router *Router) http.Handler {
	ctx := app.pool.Get().(*DefaultCtx)
	defer app.pool.Put(ctx)

	ctx.SetMetadata(router.Metadata...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx.SetCtx(w, r)
		if router.interceptor != nil {
			ctx.SetCallHandler(router.interceptor(ctx))
		}
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
				app.errorHandler(err, ctx)
			}
		}()
		err = router.Handler(ctx)
		if err != nil {
			app.errorHandler(err, ctx)
		}
	})
}

// UploadedFile returns the first uploaded file in the request, or nil if no
// files were uploaded.
func (ctx *DefaultCtx) UploadedFile() *storage.File {
	uploadedFile, ok := ctx.Get(FILE).(*storage.File)
	if uploadedFile == nil || !ok {
		return nil
	}
	return uploadedFile
}

// UploadedFiles returns a slice of all uploaded files in the request, or nil
// if no files were uploaded.
func (ctx *DefaultCtx) UploadedFiles() []*storage.File {
	uploadedFiles, ok := ctx.Get(FILES).([]*storage.File)
	if uploadedFiles == nil || !ok {
		return nil
	}
	return uploadedFiles
}

// UploadedFieldFile returns a map of uploaded files in the request, where the key
// is the field name, and the value is a slice of uploaded files for that field.
// If no files were uploaded, it returns nil.
func (ctx *DefaultCtx) UploadedFieldFile() map[string][]*storage.File {
	uploadedFieldFile, ok := ctx.Get(FIELD_FILES).(map[string][]*storage.File)
	if uploadedFieldFile == nil || !ok {
		return nil
	}
	return uploadedFieldFile
}

func (ctx *DefaultCtx) Redirect(uri string) error {
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		var scheme string
		if ctx.Req().TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}
		uri = scheme + ctx.Req().Host + ctx.Req().URL.String() + uri
	}
	fullUrl, err := url.Parse(uri)
	if err != nil {
		return err
	}
	http.Redirect(ctx.Res(), ctx.Req(), fullUrl.String(), http.StatusFound)
	return nil
}

// Ref returns the value associated with the given key from the request context.
//
// If no value is associated with the key, it returns nil.
func (ctx *DefaultCtx) Ref(name Provide) interface{} {
	return ctx.Get(name)
}
