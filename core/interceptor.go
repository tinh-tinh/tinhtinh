package core

import (
	"errors"

	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/middleware/storage"
)

const (
	FILE        CtxKey = "FILE"
	FILES       CtxKey = "FILES"
	FIELD_FILES CtxKey = "FIELD_FILES"
)

func FileInterceptor(opt storage.UploadFileOption) Middleware {
	return func(ctx Ctx) error {
		files, err := storage.HandlerFile(ctx.Req(), opt)
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}
		if len(files) == 0 {
			common.BadRequestException(ctx.Res(), "no file uploaded")
			return errors.New("no file uploaded")
		}

		ctx.Set(FILE, files[0])
		return ctx.Next()
	}
}

func FilesInterceptor(opt storage.UploadFileOption) Middleware {
	return func(ctx Ctx) error {
		files, err := storage.HandlerFile(ctx.Req(), opt)
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}
		if len(files) == 0 {
			return errors.New("no file uploaded")
		}

		ctx.Set(FILES, files)
		return ctx.Next()
	}
}

func FileFieldsInterceptor(opt storage.UploadFileOption, fieldFiles ...storage.FieldFile) Middleware {
	return func(ctx Ctx) error {
		files, err := storage.HandlerFile(ctx.Req(), opt, fieldFiles...)
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}
		if len(files) == 0 {
			common.BadRequestException(ctx.Res(), "no file upload")
			return errors.New("no file uploaded")
		}
		mapFiles := make(map[string][]*storage.File)
		for _, file := range files {
			if mapFiles[file.FieldName] == nil {
				mapFiles[file.FieldName] = []*storage.File{file}
			} else {
				mapFiles[file.FileName] = append(mapFiles[file.FileName], file)
			}
		}
		ctx.Set(FIELD_FILES, mapFiles)
		return ctx.Next()
	}
}
