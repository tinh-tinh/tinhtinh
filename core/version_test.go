package core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type Response struct {
	Data interface{} `json:"data"`
}

func AppVersionModule() ModuleParam {
	appController1 := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test").Version("1")

		ctrl.Get("/", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "1",
			})
		})
		return ctrl
	}

	appController2 := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test").Version("2")

		ctrl.Get("/", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "2",
			})
		})
		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController1, appController2},
		})

		return appModule
	}

	return module
}

func Test_VersionURI(t *testing.T) {
	app := CreateFactory(AppVersionModule(), "api")
	app.EnableVersioning(VersionOptions{
		Type: URIVersion,
	})

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/v1")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "1", res.Data)

	resp2, err := testClient.Get(testServer.URL + "/api/test/v2")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	json.Unmarshal(data2, &res2)
	require.Equal(t, "2", res2.Data)
}

func Test_VersionHeader(t *testing.T) {
	app := CreateFactory(AppVersionModule(), "api")
	app.EnableVersioning(VersionOptions{
		Type:   HeaderVersion,
		Header: "X-Version",
	})

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test/", nil)
	require.Nil(t, err)
	req.Header.Set("X-Version", "1")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "1", res.Data)

	req.Header.Set("X-Version", "2")
	resp2, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	json.Unmarshal(data2, &res2)
	require.Equal(t, "2", res2.Data)
}

func Test_VersionMedia(t *testing.T) {
	app := CreateFactory(AppVersionModule(), "api")
	app.EnableVersioning(VersionOptions{
		Type: MediaTypeVersion,
		Key:  "v=",
	})

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test/", nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json; v=1")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "1", res.Data)

	req.Header.Set("Accept", "application/json; v=2")
	resp2, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	json.Unmarshal(data2, &res2)
	require.Equal(t, "2", res2.Data)
}

func Test_VersionCustom(t *testing.T) {
	app := CreateFactory(AppVersionModule(), "api")
	app.EnableVersioning(VersionOptions{
		Type: CustomVersion,
		Extractor: func(r *http.Request) string {
			return r.URL.Query().Get("version")
		},
	})

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?version=1")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "1", res.Data)

	resp2, err := testClient.Get(testServer.URL + "/api/test?version=2")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	json.Unmarshal(data2, &res2)
	require.Equal(t, "2", res2.Data)
}
