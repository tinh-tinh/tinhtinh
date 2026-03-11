package core

import (
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

const (
	FILE        CtxKey = "FILE"
	FILES       CtxKey = "FILES"
	FIELD_FILES CtxKey = "FIELD_FILES"
)

// FileInterceptor is a middleware that intercepts a single file upload from the request,
// stores it to the configured storage, and sets the ctx key to FILE.
//
// If the request does not contain a file, it returns a 400 status code with the error message "no file uploaded".
//
// If the upload fails, it returns a 400 status code with the error message.
//
// If the upload is successful, it sets the ctx key to FILE and calls the next middleware in the chain.
func FileInterceptor(opt storage.UploadFileOption) Middleware {
	return func(ctx Ctx) error {
		file, err := storage.HandleFile(ctx.Req(), opt)
		if err != nil {
			return common.BadRequestException(ctx.Res(), err.Error())
		}
		if file == nil {
			return common.BadRequestException(ctx.Res(), "no file uploaded")
		}

		ctx.Set(FILE, file)
		return ctx.Next()
	}
}

// FilesInterceptor is a middleware that intercepts multiple file uploads from the request,
// stores them to the configured storage, and sets the ctx key to FILES.
//
// If the request does not contain any files, it returns a 400 status code with the error message "no file uploaded".
//
// If the upload fails, it returns a 400 status code with the error message.
//
// If the upload is successful, it sets the ctx key to FILES and calls the next middleware in the chain.
func FilesInterceptor(opt storage.UploadFileOption) Middleware {
	return func(ctx Ctx) error {
		files, err := storage.HandleFiles(ctx.Req(), opt)
		if err != nil {
			return common.BadRequestException(ctx.Res(), err.Error())
		}
		if len(files) == 0 {
			return common.BadRequestException(ctx.Res(), "no file uploaded")
		}

		ctx.Set(FILES, files)
		return ctx.Next()
	}
}

// FileFieldsInterceptor is a middleware that intercepts multiple file uploads from the request,
// stores them to the configured storage, and sets the ctx key to FIELD_FILES.
//
// If the request does not contain any files, it returns a 400 status code with the error message "no file uploaded".
//
// If the upload fails, it returns a 400 status code with the error message.
//
// If the upload is successful, it sets the ctx key to FIELD_FILES and calls the next middleware in the chain.
//
// The fieldFiles parameter is a list of FieldFile, which defines the name and maxCount of the field.
// The middleware will validate the number of files uploaded for each field and returns a 400 status code with the error message
// if the number of files exceeds the maxCount.
func FileFieldsInterceptor(opt storage.UploadFileOption, fieldFiles ...storage.FieldFile) Middleware {
	return func(ctx Ctx) error {
		files, err := storage.HandleFieldFiles(ctx.Req(), opt, fieldFiles...)
		if err != nil {
			return common.BadRequestException(ctx.Res(), err.Error())
		}
		if len(files) == 0 {
			return common.BadRequestException(ctx.Res(), "no file uploaded")
		}
		ctx.Set(FIELD_FILES, files)
		return ctx.Next()
	}
}
