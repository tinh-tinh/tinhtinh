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

type CtxKey string

const (
	Input CtxKey = "input"
)

func Body[M any](transform ...bool) Middleware {
	trans := false
	if len(transform) > 0 {
		trans = transform[0]
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload M
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				BadRequestException(w, err.Error())
				return
			}

			err = validator.Scanner(&payload, trans)
			if err != nil {
				BadRequestException(w, err.Error())
				return
			}
			ctx := context.WithValue(r.Context(), Input, payload)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func Query[M any](transform ...bool) Middleware {
	trans := false
	if len(transform) > 0 {
		trans = transform[0]
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var query M
			scanQuery(r, &query)

			fmt.Print(query)
			err := validator.Scanner(&query, trans)
			if err != nil {
				BadRequestException(w, err.Error())
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
