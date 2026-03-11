package storage_test

import (
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

func TestStoreFile(t *testing.T) {
	r := createMultipartRequest(t, map[string][]byte{
		"file": []byte("This is a test file"),
	})

	fileHeader, err := getFileHeaderFromRequest(r, "file")
	if err != nil {
		t.Fatalf("Failed to get file header: %v", err)
	}

	opt := storage.UploadFileOption{
		Storage: &storage.DiskOptions{
			Destination: func(r *http.Request, file *multipart.FileHeader) string {
				return "./uploads"
			},
			FileName: func(r *http.Request, file *multipart.FileHeader) string {
				return "testfile.txt"
			},
		},
	}

	uploadedFile, err := storage.StoreFile("file", fileHeader, r, opt)
	require.Nil(t, err)

	require.Equal(t, "testfile.txt", uploadedFile.FileName)
	require.Equal(t, "./uploads", uploadedFile.Destination)

	os.RemoveAll("./uploads")
}
