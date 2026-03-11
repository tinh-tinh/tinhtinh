package storage_test

import (
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

func TestHandleFieldFiles(t *testing.T) {
	r := createMultipartRequest(t, map[string][]byte{
		"file1": []byte("This is the first test file"),
		"file2": []byte("This is the second test file"),
	})

	opt := storage.UploadFileOption{
		Storage: &storage.DiskOptions{
			Destination: func(r *http.Request, file *multipart.FileHeader) string {
				return "./field_files"
			},
			FileName: func(r *http.Request, file *multipart.FileHeader) string {
				return file.Filename
			},
		},
		FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
			return file.Size != 0
		},
	}

	uploadedFiles, err := storage.HandleFieldFiles(r, opt, storage.FieldFile{
		Name:     "file1",
		MaxCount: 2,
	}, storage.FieldFile{
		Name:     "file2",
		MaxCount: 2,
	})
	require.Nil(t, err)

	require.Equal(t, 2, len(uploadedFiles))
	require.NotNil(t, uploadedFiles["file1"])
	require.Len(t, uploadedFiles["file1"], 1)
	require.Equal(t, "file1.txt", uploadedFiles["file1"][0].FileName)
	require.NotNil(t, uploadedFiles["file2"])
	require.Len(t, uploadedFiles["file2"], 1)
	require.Equal(t, "file2.txt", uploadedFiles["file2"][0].FileName)

	r = createMultipartRequest(t, map[string][]byte{})

	_, err = storage.HandleFieldFiles(r, opt, storage.FieldFile{
		Name:     "file",
		MaxCount: 2,
	})
	require.NotNil(t, err)

	r = createMultipartRequest(t, map[string][]byte{
		"file": {},
	})

	_, err = storage.HandleFieldFiles(r, opt, storage.FieldFile{
		Name:     "file",
		MaxCount: 2,
	})
	require.NotNil(t, err)

	os.RemoveAll("./field_files")
}
