package core

import (
	"context"
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

type UploadedFileInfo struct {
	FilePath     string
	OriginalName string
	FileSize     int64
	FileType     string
}

const DefaultDestFolder CtxKey = "uploads"

func FileInterceptor(fileOptions UploadedFileOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if fileOptions.DestFolder == "" {
				common.BadRequestException(w, "file destination folder cannot be empty")
				return
			}

			err := r.ParseMultipartForm(fileOptions.MaxSize)
			if err != nil {
				common.BadRequestException(w, "error parsing form: "+err.Error())
				return
			}
			fileHeaders := r.MultipartForm.File[fileOptions.FieldName]
			if len(fileHeaders) == 0 {
				common.BadRequestException(w, "no file provided for field "+fileOptions.FieldName)
				return
			}

			fileHeader := fileHeaders[0]

			file, err := fileHeader.Open()
			if err != nil {
				common.BadRequestException(w, "error opening file: "+err.Error())
				return
			}
			defer file.Close()

			fileHeaderBytes := make([]byte, 512)
			file.Read(fileHeaderBytes)
			fileType := http.DetectContentType(fileHeaderBytes)
			file.Seek(0, io.SeekStart)

			if len(fileOptions.AllowedMimes) > 0 {
				allowed := false
				for _, mime := range fileOptions.AllowedMimes {
					if fileType == mime {
						allowed = true
						break
					}
				}
				if !allowed {
					common.BadRequestException(w, "file type "+fileType+" is not allowed")
					return
				}
			}

			if _, err := os.Stat(fileOptions.DestFolder); os.IsNotExist(err) {
				os.MkdirAll(fileOptions.DestFolder, os.ModePerm)
			}

			destPath := filepath.Join(fileOptions.DestFolder, fileHeader.Filename)
			destFile, err := os.Create(destPath)
			if err != nil {
				common.BadRequestException(w, "error saving file: "+err.Error())
				return
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, file)
			if err != nil {
				common.BadRequestException(w, "error copying file: "+err.Error())
				return
			}

			filename := fileHeader.Filename
			if fileOptions.OriginalName != "" {
				filename = fileOptions.OriginalName
			}

			uploadedFile := &UploadedFileInfo{
				FilePath:     destPath,
				OriginalName: filename,
				FileSize:     fileHeader.Size,
				FileType:     fileType,
			}
			ctx := context.WithValue(r.Context(), DefaultDestFolder, uploadedFile)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func FilesInterceptor(fileOptions UploadedFilesOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if fileOptions.DestFolder == "" {
				common.BadRequestException(w, "file destination folder cannot be empty")
				return
			}

			err := r.ParseMultipartForm(fileOptions.MaxSize)
			if err != nil {
				common.BadRequestException(w, "error parsing form: "+err.Error())
				return
			}

			if len(fileOptions.AllowedMimes) == 0 {
				fileOptions.AllowedMimes = []string{
					"image/jpeg", "image/png",
					"video/mp4", "video/x-msvideo",
					"video/x-matroska", "video/quicktime",
					"video/x-flv",
				}
			}

			uploadedFiles := make(map[string][]UploadedFileInfo)

			for _, field := range fileOptions.FileFields {
				files := r.MultipartForm.File[field.Name]

				if len(files) > field.MaxCount {
					common.BadRequestException(w, "too many files for field "+field.Name)
					return
				}

				for _, fileHeader := range files {
					file, err := fileHeader.Open()
					if err != nil {
						common.BadRequestException(w, "error opening file: "+err.Error())
						return
					}
					defer file.Close()

					fileHeaderBytes := make([]byte, 512)
					file.Read(fileHeaderBytes)
					fileType := http.DetectContentType(fileHeaderBytes)
					file.Seek(0, io.SeekStart)

					allowed := false
					for _, mime := range fileOptions.AllowedMimes {
						if fileType == mime {
							allowed = true
							break
						}
					}
					if !allowed {
						common.BadRequestException(w, "file type "+fileType+" is not allowed")
						return
					}

					if _, err := os.Stat(fileOptions.DestFolder); os.IsNotExist(err) {
						os.MkdirAll(fileOptions.DestFolder, os.ModePerm)
					}

					destPath := filepath.Join(fileOptions.DestFolder, fileHeader.Filename)
					destFile, err := os.Create(destPath)
					if err != nil {
						common.BadRequestException(w, "error saving file: "+err.Error())
						return
					}
					defer destFile.Close()

					_, err = io.Copy(destFile, file)
					if err != nil {
						common.BadRequestException(w, "error copying file: "+err.Error())
						return
					}

					filename := fileHeader.Filename
					if fileOptions.OriginalName != "" {
						filename = fileOptions.OriginalName
					}

					uploadedFiles[field.Name] = append(uploadedFiles[field.Name], UploadedFileInfo{
						FilePath:     destPath,
						OriginalName: filename,
						FileSize:     fileHeader.Size,
						FileType:     fileType,
					})
				}
			}

			ctx := context.WithValue(r.Context(), DefaultDestFolder, uploadedFiles)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
