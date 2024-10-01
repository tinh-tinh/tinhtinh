package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tinh-tinh/tinhtinh/common"
)

type UploadedFileOptions struct {
	FieldName    string
	AllowedMimes []string
	MaxSize      int64
	DestFolder   string
	OriginalName string
}

type UploadedFilesOptions struct {
	FileFields   []FileField
	AllowedMimes []string
	MaxSize      int64
	DestFolder   string
	OriginalName string
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
		err := ctx.Req().ParseMultipartForm(opt.MaxSize)
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}

		fileHeaders := ctx.Req().MultipartForm.File[opt.FieldName]
		if len(fileHeaders) == 0 {
			common.BadRequestException(ctx.Res(), "no file provide for field")
			return errors.New("no file provide for field")
		}

		fileHeader := fileHeaders[0]

		file, err := fileHeader.Open()
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}
		defer file.Close()

		fileHeaderBytes := make([]byte, 512)
		_, err = file.Read(fileHeaderBytes)
		if err != nil {
			common.InternalServerException(ctx.Res(), err.Error())
			return err
		}
		fileType := http.DetectContentType(fileHeaderBytes)
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			common.InternalServerException(ctx.Res(), err.Error())
			return err
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
				return errors.New("file type " + fileType + " is not allowed")
			}
		}

		if _, err := os.Stat(opt.DestFolder); os.IsNotExist(err) {
			err = os.MkdirAll(opt.DestFolder, os.ModePerm)
			if err != nil {
				common.InternalServerException(ctx.Res(), err.Error())
				return err
			}
		}

		destPath := filepath.Join(opt.DestFolder, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}
		defer destFile.Close()

		uploadFile := &FileInfo{
			FileName:     fileHeader.Filename,
			OriginalName: opt.OriginalName,
			Encoding:     fileHeader.Header.Get("Content-Encoding"),
			MimeType:     fileType,
			Size:         fileHeader.Size,
			Destination:  destPath,
			FieldName:    opt.FieldName,
			Path:         destPath,
		}

		ctx.Set(FILE, uploadFile)
		return ctx.Next()
	}
}

func FilesInterceptor(fileOptions UploadedFilesOptions) Middleware {
	return func(ctx Ctx) error {
		err := ctx.Req().ParseMultipartForm(fileOptions.MaxSize)
		if err != nil {
			common.BadRequestException(ctx.Res(), err.Error())
			return err
		}

		if len(fileOptions.AllowedMimes) == 0 {
			fileOptions.AllowedMimes = []string{
				"image/jpeg", "image/png",
				"video/mp4", "video/x-msvideo",
				"video/x-matroska", "video/quicktime",
				"video/x-flv",
			}
		}

		uploadedFiles := make(map[string][]*FileInfo)

		for _, field := range fileOptions.FileFields {
			fileHeaders := ctx.Req().MultipartForm.File[field.Name]
			if len(fileHeaders) == 0 {
				common.BadRequestException(ctx.Res(), "no file provide for field")
				return errors.New("no file provide for field")
			}
			for i, fileHeader := range fileHeaders {
				file, err := fileHeader.Open()
				if err != nil {

					common.BadRequestException(ctx.Res(), err.Error())
					return err
				}

				fileHeaderBytes := make([]byte, 512)
				_, err = file.Read(fileHeaderBytes)
				if err != nil {
					common.InternalServerException(ctx.Res(), err.Error())
					return err
				}
				fileType := http.DetectContentType(fileHeaderBytes)
				_, err = file.Seek(0, io.SeekStart)
				if err != nil {
					common.InternalServerException(ctx.Res(), err.Error())
					return err
				}

				if len(fileOptions.AllowedMimes) > 0 {
					allowed := false
					for _, mime := range fileOptions.AllowedMimes {
						if fileType == mime {
							allowed = true
							break
						}
					}
					if !allowed {
						common.BadRequestException(ctx.Res(), "file type "+fileType+" is not allowed")
						return errors.New("file type " + fileType + " is not allowed")
					}
				}
				destPath := filepath.Join(fileOptions.DestFolder, fmt.Sprint(i), fileHeader.Filename)
				destFile, err := os.Create(destPath)
				if err != nil {
					common.BadRequestException(ctx.Res(), err.Error())
					return err
				}

				uploadFile := &FileInfo{
					FileName:     fileHeader.Filename,
					OriginalName: fileOptions.OriginalName,
					Encoding:     fileHeader.Header.Get("Content-Encoding"),
					MimeType:     fileType,
					Size:         fileHeader.Size,
					Destination:  destPath,
					FieldName:    field.Name,
					Path:         destPath,
				}

				uploadedFiles[field.Name] = append(uploadedFiles[field.Name], uploadFile)
				destFile.Close()
				file.Close()
			}
		}

		ctx.Set(FILES, uploadedFiles)
		return ctx.Next()
	}
}
