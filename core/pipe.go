package core

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

type CtxKey string

const (
	InBody  CtxKey = "body"
	InQuery CtxKey = "query"
	InPath  CtxKey = "path"
)

type Pipe[P any] struct {
	// In is the CtxKey of the Pipe. It can be one of InBody, InQuery, InPath.
	In CtxKey
}

type PipeDto interface {
	GetLocation() CtxKey
	GetValue() interface{}
}

// PipeMiddleware returns a middleware that parse and validate the given Pipes.
// The middleware will parse the request body, query or path into the given Dto
// based on the In field of the Pipe. After parsing, the middleware will validate
// the Dto using the validator.Scanner. If the validation failed, it will return
// a 400 status code with the error message. If the validation is successful, it
// will set the Dto to the ctx with the key being the In field of the Pipe and
// call the next middleware in the chain.
func PipeMiddleware(pipes ...PipeDto) Middleware {
	return func(ctx Ctx) error {
		for _, pipe := range pipes {
			dto := pipe.GetValue()
			// Clear old value in dto
			// p := reflect.ValueOf(dto).Elem()
			// p.Set(reflect.Zero(p.Type()))
			switch pipe.GetLocation() {
			case InBody:
				err := ctx.BodyParser(dto)
				if err != nil {
					return common.BadRequestException(ctx.Res(), err.Error())
				}
			case InQuery:
				err := ctx.QueryParser(dto)
				if err != nil {
					return common.BadRequestException(ctx.Res(), err.Error())
				}
			case InPath:
				err := ctx.PathParser(dto)
				if err != nil {
					return common.BadRequestException(ctx.Res(), err.Error())
				}
			}

			err := validator.Scanner(dto)
			if err != nil {
				return common.BadRequestException(ctx.Res(), err.Error())
			}
			ctx.Set(pipe.GetLocation(), dto)
		}
		return ctx.Next()
	}
}

func (p Pipe[P]) GetLocation() CtxKey {
	return p.In
}

func (p Pipe[P]) GetValue() interface{} {
	var payload P
	return &payload
}

// Body returns a Pipe that parses the request body into the given Dto.
//
// If the request body is empty, it will return a 400 status code with the error message "empty request body".
//
// If the parsing fails, it will return a 400 status code with the error message.
//
// If the parsing is successful, it will set the Dto to the ctx with the key being InBody and call the next middleware in the chain.
func Body[P any](dto P) PipeDto {
	return &Pipe[P]{
		In: InBody,
	}
}

// Query returns a Pipe that parses the query string into the given Dto.
//
// If the parsing fails, it will return a 400 status code with the error message.
//
// If the parsing is successful, it will set the Dto to the ctx with the key being InQuery and call the next middleware in the chain.
func Query[P any](dto P) PipeDto {
	return Pipe[P]{
		In: InQuery,
	}
}

// Param returns a Pipe that parses the URL path parameters into the given Dto.
//
// If the parsing fails, it will return a 400 status code with the error message.
//
// If the parsing is successful, it will set the Dto to the ctx with the key being InPath and call the next middleware in the chain.
func Param[P any](dto P) PipeDto {
	return &Pipe[P]{
		In: InPath,
	}
}

// Body Parser
type BodyParser[P any] struct{}

func (b BodyParser[P]) GetValue() any {
	var payload P
	return &payload
}

func (b BodyParser[P]) GetLocation() CtxKey {
	return InBody
}

// Query Parser
type QueryParser[P any] struct{}

func (b QueryParser[P]) GetValue() any {
	var payload P
	return &payload
}

func (b QueryParser[P]) GetLocation() CtxKey {
	return InQuery
}

// Path Parser
type PathParser[P any] struct{}

func (b PathParser[P]) GetValue() any {
	var payload P
	return &payload
}

func (b PathParser[P]) GetLocation() CtxKey {
	return InPath
}

func bindSingle(val string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)

	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Int:
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))

	case reflect.Int8:
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return err
		}
		field.SetInt(i)

	case reflect.Int16:
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return err
		}
		field.SetInt(i)

	case reflect.Int32:
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(i)

	case reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)

	case reflect.Uint:
		u, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return err
		}
		field.SetUint(u)

	case reflect.Uint8:
		u, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return err
		}
		field.SetUint(u)

	case reflect.Uint16:
		u, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return err
		}
		field.SetUint(u)

	case reflect.Uint32:
		u, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(u)

	case reflect.Uint64:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)

	case reflect.Float32:
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		}
		field.SetFloat(float64(float32(f))) // explicitly cast to float32 before setting

	case reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

func bindSlice(values []string, field reflect.Value) error {
	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(field.Type(), 0, len(values))

	for _, val := range values {
		elem := reflect.New(elemType).Elem()
		if err := bindSingle(val, elem); err != nil {
			return err
		}
		slice = reflect.Append(slice, elem)
	}

	field.Set(slice)
	return nil
}
