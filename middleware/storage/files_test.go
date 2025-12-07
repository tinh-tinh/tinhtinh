package storage_test

import (
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

func TestHandleFiles(t *testing.T) {
	r := createMultipartRequest(t, map[string][]byte{
		"file": []byte("This is the first test file"),
	})

	opt := storage.UploadFileOption{
		Storage: &storage.DiskOptions{
			Destination: func(r *http.Request, file *multipart.FileHeader) string {
				return "./files"
			},
			FileName: func(r *http.Request, file *multipart.FileHeader) string {
				return file.Filename
			},
		},
		FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
			return file.Size != 0
		},
	}

	uploadedFiles, err := storage.HandleFiles(r, opt)
	require.Nil(t, err)

	require.Equal(t, 1, len(uploadedFiles))
	require.Equal(t, "file.txt", uploadedFiles[0].FileName)

	r = createMultipartRequest(t, map[string][]byte{
		"invalid": []byte("This is an invalid test file"),
	})

	_, err = storage.HandleFiles(r, opt)
	require.NotNil(t, err)

	r = createMultipartRequest(t, map[string][]byte{
		"file": {},
	})

	_, err = storage.HandleFiles(r, opt)
	require.NotNil(t, err)

	os.RemoveAll("./files")
}
