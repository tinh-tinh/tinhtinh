package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	// Max of each file in upload
	FieldSize int64
	// Number of fields
	Fields int
	// Max of file
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

func StoreFile(field string, fileHeader *multipart.FileHeader, r *http.Request, opt UploadFileOption) (*File, error) {
	var destFolder string
	if opt.Storage.Destination != nil {
		destFolder = opt.Storage.Destination(r, fileHeader)
	}

	if destFolder != "" {
		if _, err := os.Stat(destFolder); os.IsNotExist(err) {
			err := os.MkdirAll(destFolder, 0o755)
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

	// Open the source file and copy its content to the destination
	srcFile, err := fileHeader.Open()
	if err != nil {
		// Clean up the empty destination file
		os.Remove(destPath)
		return nil, err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		// Clean up the incomplete destination file
		os.Remove(destPath)
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	return &File{
		FieldName:    field,
		FileName:     fileName,
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
		Destination:  destFolder,
		Path:         destPath,
	}, nil
}
