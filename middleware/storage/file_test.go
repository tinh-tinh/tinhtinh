package storage_test

import (
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

func TestHandleFile(t *testing.T) {
	r := createMultipartRequest(t, map[string][]byte{
		"file": []byte("This is a test file"),
	})

	opt := storage.UploadFileOption{
		Storage: &storage.DiskOptions{
			Destination: func(r *http.Request, file *multipart.FileHeader) string {
				return "./file"
			},
			FileName: func(r *http.Request, file *multipart.FileHeader) string {
				return "testfile.txt"
			},
		},
		FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
			return strings.HasPrefix(file.Filename, "file")
		},
	}

	uploadedFile, err := storage.HandleFile(r, opt)
	require.Nil(t, err)

	require.Equal(t, "testfile.txt", uploadedFile.FileName)
	require.Equal(t, "./file", uploadedFile.Destination)

	r = createMultipartRequest(t, map[string][]byte{
		"invalid": []byte("This is a test file"),
	})

	_, err = storage.HandleFile(r, opt)
	require.NotNil(t, err)

	os.RemoveAll("./file")
}
