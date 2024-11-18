package core

import (
	"reflect"

	"github.com/tinh-tinh/tinhtinh/common"
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
	// Dto is the data that will be parsed and validated.
	Dto interface{}
	// In is the InDto of the Pipe. It can be one of InBody, InQuery, InPath.
	In InDto
}

// PipeMiddleware returns a middleware that parse and validate the given Pipes.
// The middleware will parse the request body, query or path into the given Dto
// based on the In field of the Pipe. After parsing, the middleware will validate
// the Dto using the validator.Scanner. If the validation failed, it will return
// a 400 status code with the error message. If the validation is successful, it
// will set the Dto to the ctx with the key being the In field of the Pipe and
// call the next middleware in the chain.
func PipeMiddleware(pipes ...Pipe) Middleware {
	return func(ctx Ctx) error {
		for _, pipe := range pipes {
			dto := pipe.Dto
			// Clear old value in dto
			p := reflect.ValueOf(dto).Elem()
			p.Set(reflect.Zero(p.Type()))
			switch pipe.In {
			case InBody:
				err := ctx.BodyParser(dto)
				if err != nil {
					return common.BadRequestException(ctx.Res(), err.Error())
				}
			case InQuery:
				err := ctx.QueryParse(dto)
				if err != nil {
					return common.BadRequestException(ctx.Res(), err.Error())
				}
			case InPath:
				err := ctx.ParamParse(dto)
				if err != nil {
					return common.BadRequestException(ctx.Res(), err.Error())
				}
			}

			err := validator.Scanner(dto)
			if err != nil {
				return common.BadRequestException(ctx.Res(), err.Error())
			}
			ctx.Set(pipe.In, dto)
		}
		return ctx.Next()
	}
}

// Body returns a Pipe that parses the request body into the given Dto.
//
// If the request body is empty, it will return a 400 status code with the error message "empty request body".
//
// If the parsing fails, it will return a 400 status code with the error message.
//
// If the parsing is successful, it will set the Dto to the ctx with the key being InBody and call the next middleware in the chain.
func Body(dto interface{}) Pipe {
	return Pipe{
		Dto: dto,
		In:  InBody,
	}
}

// Query returns a Pipe that parses the query string into the given Dto.
//
// If the parsing fails, it will return a 400 status code with the error message.
//
// If the parsing is successful, it will set the Dto to the ctx with the key being InQuery and call the next middleware in the chain.
func Query(dto interface{}) Pipe {
	return Pipe{
		Dto: dto,
		In:  InQuery,
	}
}

// Param returns a Pipe that parses the URL path parameters into the given Dto.
//
// If the parsing fails, it will return a 400 status code with the error message.
//
// If the parsing is successful, it will set the Dto to the ctx with the key being InPath and call the next middleware in the chain.
func Param(dto interface{}) Pipe {
	return Pipe{
		Dto: dto,
		In:  InPath,
	}
}
