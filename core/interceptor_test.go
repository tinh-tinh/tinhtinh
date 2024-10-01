package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FileInterceptor(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(FileInterceptor(UploadedFileOptions{
			DestFolder: "./upload",
			MaxSize:    1024,
			FieldName:  "file",
		})).Post("", func(ctx Ctx) {

			fmt.Println(ctx.UploadedFile())
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

	req, err := http.NewRequest("POST", testServer.URL+"/api/test", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	files, er := os.Open("./upload")
	require.Nil(t, er)

	defer files.Close()

	filesInfo, err := files.Readdir(-1)
	require.Nil(t, err)

	for _, file := range filesInfo {
		os.Remove(filepath.Join("./upload", file.Name()))
	}
}
