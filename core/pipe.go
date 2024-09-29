package core

import (
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
	Dto interface{}
	In  InDto
}

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
				err := ctx.QueryParse(dto)
				if err != nil {
					common.BadRequestException(ctx.Res(), err.Error())
					return err
				}
			case InPath:
				err := ctx.ParamParse(dto)
				if err != nil {
					common.BadRequestException(ctx.Res(), err.Error())
					return err
				}
			}

			err := validator.Scanner(dto)
			if err != nil {
				common.BadRequestException(ctx.Res(), err.Error())
				return err
			}
			ctx.Set(pipe.In, dto)
		}
		return ctx.Next()
	}
}

func Body(dto interface{}) Pipe {
	return Pipe{
		Dto: dto,
		In:  InBody,
	}
}

func Query(dto interface{}) Pipe {
	return Pipe{
		Dto: dto,
		In:  InQuery,
	}
}

func Param(dto interface{}) Pipe {
	return Pipe{
		Dto: dto,
		In:  InPath,
	}
}
