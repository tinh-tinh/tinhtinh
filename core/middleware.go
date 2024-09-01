package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/dto/validator"
	"github.com/tinh-tinh/tinhtinh/utils"
)

type Middleware func(http.Handler) http.Handler

type Interceptor func(ctx Ctx) http.Handler

type CtxKey string

const (
	Input CtxKey = "input"
)

func Body[M any](transform ...bool) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload M
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = validator.Scanner(&payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), Input, payload)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func Query[M any]() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var query M
			scanQuery(r, &query)

			fmt.Print(query)
			err := validator.Scanner(&query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func scanQuery(r *http.Request, payload interface{}) {
	ct := reflect.ValueOf(payload).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("query")
		if tagVal != "" {
			val := r.URL.Query().Get(tagVal)
			if val == "" {
				continue
			}
			switch field.Type.Name() {
			case "string":
				ct.Field(i).SetString(val)
			case "int":
				ct.Field(i).SetInt(int64(utils.StringToInt(val)))
			case "bool":
				ct.Field(i).SetBool(utils.StringToBool(val))
			default:
				fmt.Println(field.Type.Name())
			}
		}
	}
}

// Guard func
type Guard func(ctx Ctx) bool

func ParseGuard(guard Guard) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAccess := guard(NewCtx(w, r))
			if !isAccess {
				http.Error(w, "You can not access", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
