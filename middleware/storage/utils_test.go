package storage_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

func createMultipartRequest(t *testing.T, files map[string][]byte) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, fileContent := range files {
		part, err := writer.CreateFormFile(fieldName, fieldName+".txt")
		require.Nil(t, err)
		_, err = part.Write(fileContent)
		require.Nil(t, err)
	}

	err := writer.Close()
	require.Nil(t, err)

	req, err := http.NewRequest("POST", "/upload", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.FormFile("file")

	return req
}

func TestValidateLimit(t *testing.T) {
	req := createMultipartRequest(t, map[string][]byte{
		"file1": []byte("This is file 1"),
		"file2": []byte("This is file 2"),
	})

	// Test with no limit
	err := storage.ValidateLimit(nil, req)
	require.Nil(t, err)

	// Test with file size limit
	// reqWithBigSize := createMultipartRequest(t, map[string][]byte{
	// 	"file": bytes.Repeat([]byte("A"), 5<<20), // 5 MB
	// })
	limit := &storage.UploadFileLimit{FileSize: 1}
	// err = storage.ValidateLimit(limit, reqWithBigSize)
	// require.NotNil(t, err)

	// Test with field limit
	limit = &storage.UploadFileLimit{Fields: 1}
	err = storage.ValidateLimit(limit, req)
	require.NotNil(t, err)

	limit = &storage.UploadFileLimit{Fields: 2}
	err = storage.ValidateLimit(limit, req)
	require.Nil(t, err)
}

func getFileHeaderFromRequest(r *http.Request, fieldName string) (*multipart.FileHeader, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}

	files := r.MultipartForm.File[fieldName]
	if len(files) == 0 {
		return nil, http.ErrMissingFile
	}

	return files[0], nil
}

func TestValidateFilterFile(t *testing.T) {
	req := createMultipartRequest(t, map[string][]byte{
		"file": []byte("This is a test file"),
	})

	fileHeader, err := getFileHeaderFromRequest(req, "file")
	require.Nil(t, err)

	// Test with no filter
	opt := storage.UploadFileOption{}
	err = storage.ValidateFilterFile(req, fileHeader, opt)
	require.Nil(t, err)

	// Test with passing filter
	opt = storage.UploadFileOption{
		FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
			return true
		},
	}
	err = storage.ValidateFilterFile(req, fileHeader, opt)
	require.Nil(t, err)

	// Test with failing filter
	opt = storage.UploadFileOption{
		FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
			return false
		},
	}
	err = storage.ValidateFilterFile(req, fileHeader, opt)
	require.NotNil(t, err)
}

func TestDetectMime(t *testing.T) {
	tests := []struct {
		filename     string
		content      []byte
		expectedMime string
	}{
		{"test.txt", []byte("This is a plain text file."), "text/plain; charset=utf-8"},
		{"image.png", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, "image/png"},
		{"unknown.bin", []byte{0x00, 0x01, 0x02, 0x03}, "application/octet-stream"},
		{"noext", []byte("%PDF-1.4"), "application/pdf"},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			rs := bytes.NewReader(test.content)
			mimeType, err := storage.DetectAndValidateContentType(rs, test.filename)
			require.Nil(t, err)
			require.Equal(t, test.expectedMime, mimeType)
		})
	}
}
