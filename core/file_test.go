package core_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
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

	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(core.FileInterceptor(storage.UploadFileOption{
			Storage: store,
		})).Post("happy", func(ctx core.Ctx) error {

			return ctx.JSON(core.Map{
				"data": ctx.UploadedFile().OriginalName,
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
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

func Test_FilesInterceptor(t *testing.T) {
	store := &storage.DiskOptions{
		Destination: func(r *http.Request, file *multipart.FileHeader) string {
			return "./upload"
		},
		FileName: func(r *http.Request, file *multipart.FileHeader) string {
			uniqueSuffix := time.Now().Format("20060102150405") + "-" + file.Filename
			return uniqueSuffix
		},
	}

	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(core.FilesInterceptor(storage.UploadFileOption{
			Storage: store,
		})).Post("", func(ctx core.Ctx) error {
			files := ctx.UploadedFiles()

			fmt.Printf("file is %v", files)
			return ctx.JSON(core.Map{
				"data": len(files),
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Prepare test file
	// Prepare the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add multiple files to the form
	tempFiles := []struct {
		fieldName string
		fileName  string
		content   string
	}{
		{"file", "test1.txt", "Hello, World!"},
		{"file", "test2.txt", "Another test file"},
	}

	for _, file := range tempFiles {
		part, err := writer.CreateFormFile(file.fieldName, file.fileName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, bytes.NewBufferString(file.content))
		if err != nil {
			t.Fatal(err)
		}
	}

	err := writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Happy case
	req, err := http.NewRequest("POST", testServer.URL+"/api/test", body)
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
	require.Equal(t, float64(2), res.Data)

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

func Test_FieldFileInterceptor(t *testing.T) {
	store := &storage.DiskOptions{
		Destination: func(r *http.Request, file *multipart.FileHeader) string {
			return "./upload"
		},
		FileName: func(r *http.Request, file *multipart.FileHeader) string {
			uniqueSuffix := time.Now().Format("20060102150405") + "-" + file.Filename
			return uniqueSuffix
		},
	}

	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(
			core.FileFieldsInterceptor(storage.UploadFileOption{
				Storage: store,
			}, storage.FieldFile{
				Name:     "file1",
				MaxCount: 2,
			}, storage.FieldFile{
				Name:     "file2",
				MaxCount: 2,
			}),
		).Post("", func(ctx core.Ctx) error {
			files := ctx.UploadedFieldFile()

			idx := 0
			for k, v := range files {
				for _, file := range v {
					fmt.Printf("Filed %s with file %v\n", k, file.FileName)
					idx++
				}
			}

			return ctx.JSON(core.Map{
				"data": idx,
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Prepare test file
	// Prepare the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add multiple files to the form
	tempFiles := []struct {
		fieldName string
		fileName  string
		content   string
	}{
		{"file1", "test1.txt", "Hello, World!"},
		{"file2", "test2.txt", "Another test file"},
	}

	for _, file := range tempFiles {
		part, err := writer.CreateFormFile(file.fieldName, file.fileName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, bytes.NewBufferString(file.content))
		if err != nil {
			t.Fatal(err)
		}
	}

	err := writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Happy case
	req, err := http.NewRequest("POST", testServer.URL+"/api/test", body)
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
	require.Equal(t, float64(2), res.Data)

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
