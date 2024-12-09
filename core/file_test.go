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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

func uploadFile(name string) (*bytes.Buffer, string, string) {
	// Add multiple files to the form
	tempFile, err := os.CreateTemp("", name)
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())

	content := []byte("This is a test file content")
	if _, err := tempFile.Write(content); err != nil {
		panic(err)
	}
	tempFile.Close()

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile.Name()))
	if err != nil {
		panic(err)
	}
	file, err := os.Open(tempFile.Name())
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	writer.Close()
	return body, writer.FormDataContentType(), filepath.Base(tempFile.Name())
}

type testFile struct {
	fieldName string
	fileName  string
	content   string
}

func uploadFiles(tempFiles []testFile) (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, file := range tempFiles {
		part, err := writer.CreateFormFile(file.fieldName, file.fileName)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(part, bytes.NewBufferString(file.content))
		if err != nil {
			panic(err)
		}
	}

	err := writer.Close()
	if err != nil {
		panic(err)
	}

	return body, writer
}

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
			FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
				return strings.HasPrefix(file.Filename, "test")
			},
		})).Post("happy", func(ctx core.Ctx) error {

			return ctx.JSON(core.Map{
				"data": ctx.UploadedFile().OriginalName,
			})
		})

		ctrl.Post("badluck", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.UploadedFile(),
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Case 1: Happy Case
	body, contentType, fileName := uploadFile("test-upload-*.txt")

	req, err := http.NewRequest("POST", testServer.URL+"/api/test/happy", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", contentType)

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, fileName, res.Data)

	// Case 2: Not use middleware
	resp, err = testClient.Post(testServer.URL+"/api/test/badluck", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":null}`, string(data))

	// Case 3: Filter Failed
	body3, contentType3, _ := uploadFile("upload-*.txt")

	req, err = http.NewRequest("POST", testServer.URL+"/api/test/happy", body3)
	require.Nil(t, err)
	req.Header.Set("Content-Type", contentType3)

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintln(`{"error":"file filter failed"}`), string(data))

	// Case 4: Not File Upload
	body4 := &bytes.Buffer{}
	writer4 := multipart.NewWriter(body4)
	writer4.Close()

	req, err = http.NewRequest("POST", testServer.URL+"/api/test/happy", body4)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer4.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintln(`{"error":"no file uploaded"}`), string(data))

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
			FileFilter: func(r *http.Request, file *multipart.FileHeader) bool {
				return strings.HasPrefix(file.Filename, "test")
			},
		})).Post("", func(ctx core.Ctx) error {
			files := ctx.UploadedFiles()

			return ctx.JSON(core.Map{
				"data": len(files),
			})
		})

		ctrl.Post("badluck", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.UploadedFiles(),
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Case 1: Happy case
	body, writer := uploadFiles([]testFile{
		{"file", "test1.txt", "Hello, World!"},
		{"file", "test2.txt", "Another test file"},
	})

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

	// Case 2: Not use middleware
	resp, err = testClient.Post(testServer.URL+"/api/test/badluck", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":null}`, string(data))

	// Case 3: Filter Failed
	body3, writer3 := uploadFiles([]testFile{
		{"file", "1.txt", "Hello, World!"},
		{"file", "2.txt", "Another test file"},
	})

	req, err = http.NewRequest("POST", testServer.URL+"/api/test", body3)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer3.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintln(`{"error":"file filter failed"}`), string(data))

	// Case 4: Not File Upload
	body4 := &bytes.Buffer{}
	writer4 := multipart.NewWriter(body4)
	writer4.Close()

	req, err = http.NewRequest("POST", testServer.URL+"/api/test", body4)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer4.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintln(`{"error":"no file uploaded"}`), string(data))

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
					fmt.Printf("Field %s with file %v\n", k, file.FileName)
					idx++
				}
			}

			return ctx.JSON(core.Map{
				"data": idx,
			})
		})

		ctrl.Post("badluck", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.UploadedFieldFile(),
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	body, writer := uploadFiles([]testFile{
		{"file1", "test1.txt", "Hello, World!"},
		{"file2", "test2.txt", "Another test file"},
		{"file1", "test2.txt", "Hello, World!"},
	})

	// Case 1: Happy case
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
	require.Equal(t, float64(3), res.Data)

	// Case 2: Not use middleware
	resp, err = testClient.Post(testServer.URL+"/api/test/badluck", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":null}`, string(data))

	// Case 3: Limited Failed
	body2, writer2 := uploadFiles([]testFile{
		{"file1", "test1.txt", "Hello, World!"},
		{"file2", "test2.txt", "Another test file"},
		{"file1", "test2.txt", "Hello, World!"},
		{"file1", "test3.txt", "Hello, World!"},
	})

	req, err = http.NewRequest("POST", testServer.URL+"/api/test", body2)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer2.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintln(`{"error":"number of field file1 exceeds limit 2"}`), string(data))

	// Case 4: Not File Upload
	body4 := &bytes.Buffer{}
	writer4 := multipart.NewWriter(body4)
	writer4.Close()

	req, err = http.NewRequest("POST", testServer.URL+"/api/test", body4)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer4.FormDataContentType())

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintln(`{"error":"no file uploaded"}`), string(data))

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
