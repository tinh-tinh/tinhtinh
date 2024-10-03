package storage

import (
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Storage struct {
	Destination string
	FileName    string
}

type Callback func(error, string)

type DiskOptions struct {
	Destination func(r *http.Request, file *multipart.FileHeader) string
	FileName    func(r *http.Request, file *multipart.FileHeader) string
}

type UploadFileLimit struct {
	// Default 1MB
	FieldSize int64
	// Infinite
	Fields int
	// Infinite
	FileSize int64
}

type FileFilter func(r *http.Request, file *multipart.FileHeader) bool

type File struct {
	FileName     string
	OriginalName string
	Encoding     string
	MimeType     string
	Size         int64
	Destination  string
	FieldName    string
	Path         string
}

type UploadFileOption struct {
	FieldName  string
	FileFilter FileFilter
	Storage    *DiskOptions
	Limit      *UploadFileLimit
}

type FieldFile struct {
	Name     string
	MaxCount int
}

func HandlerFile(r *http.Request, opt UploadFileOption, fieldFiles ...FieldFile) ([]*File, error) {
	uploadFiles := []*File{}
	if opt.FieldName == "" {
		opt.FieldName = "file"
	}

	r.FormFile(opt.FieldName)
	if opt.Storage == nil {
		return nil, errors.New("storage is required")
	}

	// Validate limit
	if opt.Limit != nil {
		if opt.Limit.FileSize > 0 {
			err := r.ParseMultipartForm(opt.Limit.FileSize)
			if err != nil {
				return nil, err
			}
		}

		if opt.Limit.Fields > 0 {
			numFields := len(r.MultipartForm.File)
			if numFields > opt.Limit.Fields {
				return nil, errors.New("number of fields exceeds limit" + strconv.Itoa(opt.Limit.Fields))
			}
		}
	}

	isUploadSingle := len(fieldFiles) == 0
	for field, files := range r.MultipartForm.File {
		if field != opt.FieldName && isUploadSingle {
			continue
		}
		if !isUploadSingle {
			matchField := slices.IndexFunc(fieldFiles, func(e FieldFile) bool {
				return e.Name == field
			})
			if len(files) > fieldFiles[matchField].MaxCount {
				return nil, errors.New("number of fields exceeds limit" + strconv.Itoa(fieldFiles[matchField].MaxCount))
			}
		}
		for _, fileHeader := range files {
			if opt.FileFilter != nil && !opt.FileFilter(r, fileHeader) {
				return nil, errors.New("file filter failed")
			}

			if opt.Limit != nil && fileHeader.Size > opt.Limit.FieldSize {
				return nil, errors.New("file size exceeds limit" + strconv.FormatInt(opt.Limit.FieldSize, 10))
			}

			mimeType := fileHeader.Header.Get("Content-Type")
			if mimeType == "" {
				mimeType = mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
			}
			mediaType, params, err := mime.ParseMediaType(mimeType)
			if err != nil {
				return nil, err
			}
			encode := params["charset"]

			// store
			var destFolder string
			if opt.Storage.Destination != nil {
				destFolder = opt.Storage.Destination(r, fileHeader)
			}

			if destFolder != "" {
				if _, err := os.Stat(destFolder); os.IsNotExist(err) {
					err := os.MkdirAll(destFolder, 0755)
					if err != nil {
						return nil, err
					}
				}
			}

			var fileName string
			if opt.Storage.FileName != nil {
				fileName = opt.Storage.FileName(r, fileHeader)
			}

			destPath := filepath.Join(destFolder, fileName)
			destFile, err := os.Create(destPath)
			if err != nil {
				return nil, err
			}
			defer destFile.Close()

			uploadFiles = append(uploadFiles, &File{
				FieldName:    field,
				FileName:     fileName,
				OriginalName: fileHeader.Filename,
				Encoding:     encode,
				MimeType:     mediaType,
				Size:         fileHeader.Size,
				Destination:  destFolder,
				Path:         destPath,
			})
		}
	}

	return uploadFiles, nil
}
