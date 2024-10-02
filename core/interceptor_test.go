package core

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/middleware/storage"
)

func Test_FileInterceptor(t *testing.T) {
	store := &storage.DiskOptions{
		Destination: func(r *http.Request, file *multipart.FileHeader) string {
			return "./upload"
		},
		FileName: func(r *http.Request, file *multipart.FileHeader) string {
			uniqueSuffix := time.Now().Format("20060102150405") + "-" + file.Filename
			return uniqueSuffix
		},
	}

	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(FileInterceptor(UploadedFileOptions{
			Storage: store,
		})).Post("happy", func(ctx Ctx) {

			ctx.JSON(Map{
				"data": ctx.UploadedFile().OriginalName,
			})
		})

		ctrl.Use(FileInterceptor(UploadedFileOptions{})).Post("no-store", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.UploadedFile(),
			})
		})

		ctrl.Use(FileInterceptor(UploadedFileOptions{
			Storage: store,
			Limit: &UploadFileLimit{
				FileSize: 1,
			},
		})).Post("limit-size", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.UploadedFile(),
			})
		})

		ctrl.Use(FileInterceptor(UploadedFileOptions{
			Storage: store,
			AllowedMimes: []string{
				"image/*",
			},
		})).Post("allow-mime", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.UploadedFile(),
			})
		})

		ctrl.Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.UploadedFile(),
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Prepare test file
	tempFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	content := []byte("This is a test file content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile.Name()))
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	// Happy case
	req, err := http.NewRequest("POST", testServer.URL+"/api/test/happy", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, filepath.Base(tempFile.Name()), res.Data)

	// Case no file
	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Case No store
	req, err = http.NewRequest("POST", testServer.URL+"/api/test/no-store", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Case limit size
	req, err = http.NewRequest("POST", testServer.URL+"/api/test/limit-size", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Remove all file after test
	files, er := os.Open("./upload")
	require.Nil(t, er)

	defer files.Close()

	filesInfo, err := files.Readdir(-1)
	require.Nil(t, err)

	for _, file := range filesInfo {
		os.Remove(filepath.Join("./upload", file.Name()))
	}
}
