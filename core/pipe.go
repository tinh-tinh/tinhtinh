package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/dto/transform"
	"github.com/tinh-tinh/tinhtinh/dto/validator"
)

type CtxKey string

type InDto string

const (
	InBody  InDto = "body"
	InQuery InDto = "query"
	InPath  InDto = "path"
)

type Pipe struct {
	Dto       interface{}
	Transform bool
	In        InDto
}

func PipeMiddleware(pipes ...Pipe) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, pipe := range pipes {
				dto := pipe.Dto
				switch pipe.In {
				case InBody:
					err := json.NewDecoder(r.Body).Decode(dto)
					if err != nil {
						BadRequestException(w, err.Error())
						return
					}
				case InQuery:
					scanQuery(r, dto)
				case InPath:
					scanParam(r, dto)
				}

				err := validator.Scanner(dto, pipe.Transform)
				if err != nil {
					BadRequestException(w, err.Error())
					return
				}
				ctx := context.WithValue(r.Context(), pipe.In, dto)
				r = r.WithContext(ctx)
				h.ServeHTTP(w, r)
			}
		})
	}
}

func Body(dto interface{}, transform ...bool) Pipe {
	trans := false
	if len(transform) > 0 {
		trans = transform[0]
	}
	return Pipe{
		Dto:       dto,
		In:        InBody,
		Transform: trans,
	}
}

func Query(dto interface{}, transform ...bool) Pipe {
	trans := false
	if len(transform) > 0 {
		trans = transform[0]
	}
	return Pipe{
		Dto:       dto,
		In:        InQuery,
		Transform: trans,
	}
}

func Param(dto interface{}, transform ...bool) Pipe {
	trans := false
	if len(transform) > 0 {
		trans = transform[0]
	}
	return Pipe{
		Dto:       dto,
		In:        InPath,
		Transform: trans,
	}
}

func scanParam(r *http.Request, payload interface{}) {
	ct := reflect.ValueOf(payload).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("param")
		if tagVal != "" {
			val := r.PathValue(tagVal)
			if val == "" {
				continue
			}
			switch field.Type.Name() {
			case "string":
				ct.Field(i).SetString(val)
			case "int":
				ct.Field(i).SetInt(transform.StringToInt64(val))
			case "bool":
				ct.Field(i).SetBool(transform.StringToBool(val))
			default:
				fmt.Println(field.Type.Name())
			}
		}
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
				ct.Field(i).SetInt(transform.StringToInt64(val))
			case "bool":
				ct.Field(i).SetBool(transform.StringToBool(val))
			default:
				fmt.Println(field.Type.Name())
			}
		}
	}
}
