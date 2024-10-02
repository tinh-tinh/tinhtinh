package core

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/middleware/storage"
)

type UploadedFileOptions struct {
	FieldName    string
	AllowedMimes []string
	Limit        *UploadFileLimit
	Storage      *storage.DiskOptions
}

type UploadFileLimit struct {
	// Default 1MB
	FieldSize int
	// Infinite
	Fields int
	// Infinite
	FileSize int64
}

type UploadedFilesOptions struct {
	FileFields []FileField
	Option     UploadedFileOptions
}
type FileField struct {
	Name     string
	MaxCount int
}

type FileInfo struct {
	FileName     string
	OriginalName string
	Encoding     string
	MimeType     string
	Size         int64
	Destination  string
	FieldName    string
	Path         string
}

const DefaultDestFolder CtxKey = "uploads"

const FILE CtxKey = "FILE"
const FILES CtxKey = "FILES"

func FileInterceptor(opt UploadedFileOptions) Middleware {

	return func(ctx Ctx) error {
		file, err := handerFile(ctx, opt)
		if err != nil {
			return err
		}

		ctx.Set(FILE, file)
		return ctx.Next()
	}
}

func FilesInterceptor(key string, opt UploadedFileOptions) Middleware {
	return func(ctx Ctx) error {
		files := []FileInfo{}

		if opt.FieldName == "" {
			opt.FieldName = "file"
		}
		if opt.Storage == nil {
			common.BadRequestException(ctx.Res(), "storage is required")
			return errors.New("storage is required")
		}

		// Validate
		if opt.Limit != nil {
			if opt.Limit.FileSize > 0 {
				err := ctx.Req().ParseMultipartForm(opt.Limit.FileSize)
				if err != nil {
					common.BadRequestException(ctx.Res(), err.Error())
					return err
				}
			}
		}

		fileUploads := ctx.Req().MultipartForm.File[opt.FieldName]
		for _, fHeader := range fileUploads {
			// Open the file
			file, err := fHeader.Open()
			if err != nil {
				// Handle error
			}
			defer file.Close()
			// Process each file similarly to the single file scenario
		}

		ctx.Set(FILES, files)
		return ctx.Next()
	}
}

func handerFile(ctx Ctx, opt UploadedFileOptions) (*FileInfo, error) {
	if opt.FieldName == "" {
		opt.FieldName = "file"
	}
	if opt.Storage == nil {
		common.BadRequestException(ctx.Res(), "storage is required")
		return nil, errors.New("storage is required")
	}

	// Validate
	if opt.Limit != nil {
		if opt.Limit.FileSize > 0 {
			err := ctx.Req().ParseMultipartForm(opt.Limit.FileSize)
			if err != nil {
				common.BadRequestException(ctx.Res(), err.Error())
				return nil, err
			}
		}
	}

	file, handler, err := ctx.Req().FormFile(opt.FieldName)
	if err != nil {
		common.InternalServerException(ctx.Res(), err.Error())
		return nil, err
	}
	defer file.Close()

	fileHeaderBytes := make([]byte, 512)
	_, err = file.Read(fileHeaderBytes)
	if err != nil {
		common.InternalServerException(ctx.Res(), err.Error())
		return nil, err
	}
	fileType := http.DetectContentType(fileHeaderBytes)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		common.InternalServerException(ctx.Res(), err.Error())
		return nil, err
	}

	if len(opt.AllowedMimes) > 0 {
		allowed := false
		for _, mime := range opt.AllowedMimes {
			if fileType == mime {
				allowed = true
				break
			}
		}
		if !allowed {
			common.BadRequestException(ctx.Res(), "file type "+fileType+" is not allowed")
			return nil, errors.New("file type " + fileType + " is not allowed")
		}
	}

	// Store
	var destFolder string
	if opt.Storage.Destination != nil {
		destFolder = opt.Storage.Destination(ctx.Req(), handler)
	}

	if destFolder != "" {
		if _, err := os.Stat(destFolder); os.IsNotExist(err) {
			err = os.MkdirAll(destFolder, os.ModePerm)
			if err != nil {
				common.InternalServerException(ctx.Res(), err.Error())
				return nil, err
			}
		}
	}

	var fileName string
	if opt.Storage.FileName != nil {
		fileName = opt.Storage.FileName(ctx.Req(), handler)
	}

	destPath := filepath.Join(destFolder, fileName)
	destFile, err := os.Create(destPath)
	if err != nil {
		common.BadRequestException(ctx.Res(), err.Error())
		return nil, err
	}
	defer destFile.Close()

	uploadFile := &FileInfo{
		FileName:     fileName,
		OriginalName: handler.Filename,
		Encoding:     handler.Header.Get("Content-Encoding"),
		MimeType:     fileType,
		Size:         handler.Size,
		Destination:  destFolder,
		FieldName:    opt.FieldName,
		Path:         destPath,
	}

	return uploadFile, nil
}
