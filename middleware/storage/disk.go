package storage

import (
	"mime/multipart"
	"net/http"
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
