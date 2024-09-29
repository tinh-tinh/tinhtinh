package core

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/common"
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

// PipeMiddleware returns a middleware that populates the given dtos from the
// request body, query or path and validates them.
//
// The middleware will set the dtos as values of the request context with the
// InDto as the key. The dtos are then accessible in the next middleware and
// handlers via the context.
//
// The middleware will also validate the dtos and return a 400 status code if
// any of the dtos are invalid.
func PipeMiddleware(pipes ...Pipe) Middleware {
	return func(ctx Ctx) error {
		for _, pipe := range pipes {
			dto := pipe.Dto
			switch pipe.In {
			case InBody:
				err := ctx.BodyParser(dto)
				if err != nil {
					common.BadRequestException(ctx.Res(), err.Error())
					return err
				}
			case InQuery:
				scanQuery(ctx.Req(), dto)
			case InPath:
				scanParam(ctx.Req(), dto)
			}

			fmt.Println(dto)
			err := validator.Scanner(dto, pipe.Transform)
			if err != nil {
				common.BadRequestException(ctx.Res(), err.Error())
				return err
			}
			ctx.Set(pipe.In, dto)
		}
		return ctx.Next()
	}
}

// Body returns a Pipe that populates the given dto from the request body and
// validates it.
//
// The dto is populated from the request body by decoding the JSON. The dto is
// then validated by the validator.Scanner. If the dto is invalid, the
// middleware will return a 400 status code and write the error message to the
// response.
//
// The transform parameter is passed to the validator.Scanner as the transform
// option. If it is true, the validator.Scanner will transform the dto according
// to the rules specified in the dto/transform package.
//
// The Pipe will set the dto as a value of the request context with the key
// InBody. The dto is then accessible in the next middleware and handlers via
// the context.
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

// Query returns a Pipe that populates the given dto from the query parameters
// and validates it.
//
// The dto is populated from the query parameters by scanning the query string.
// The dto is then validated by the validator.Scanner. If the dto is invalid, the
// middleware will return a 400 status code and write the error message to the
// response.
//
// The transform parameter is passed to the validator.Scanner as the transform
// option. If it is true, the validator.Scanner will transform the dto according
// to the rules specified in the dto/transform package.
//
// The Pipe will set the dto as a value of the request context with the key
// InQuery. The dto is then accessible in the next middleware and handlers via
// the context.
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

// Param returns a Pipe that populates the given dto from the path parameters
// and validates it.
//
// The dto is populated from the path parameters by scanning the path.
// The dto is then validated by the validator.Scanner. If the dto is invalid, the
// middleware will return a 400 status code and write the error message to the
// response.
//
// The transform parameter is passed to the validator.Scanner as the transform
// option. If it is true, the validator.Scanner will transform the dto according
// to the rules specified in the dto/transform package.
//
// The Pipe will set the dto as a value of the request context with the key
// InPath. The dto is then accessible in the next middleware and handlers via
// the context.
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

// scanParam scans the path parameters of the given request and sets the
// corresponding field in the given payload.
//
// The path parameter name is determined by the "param" tag on the field.
// If the tag is not present or the parameter is not present in the path,
// the field is left unchanged.
//
// The type of the field determines how the value is set:
//   - string: the value is set directly
//   - int: the value is parsed as an int64 and set
//   - bool: the value is parsed as a bool and set
//   - other: the value is not set
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

// scanQuery scans the query parameters of the given request and sets the
// corresponding field in the given payload.
//
// The query parameter name is determined by the "query" tag on the field.
// If the tag is not present or the parameter is not present in the query,
// the field is left unchanged.
//
// The type of the field determines how the value is set:
//   - string: the value is set directly
//   - int: the value is parsed as an int64 and set
//   - bool: the value is parsed as a bool and set
//   - other: the value is not set
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
