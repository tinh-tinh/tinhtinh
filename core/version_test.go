package core_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Response struct {
	Data interface{} `json:"data"`
}

func AppVersionModule() core.ModuleParam {
	appController1 := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test").Version("1")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "1",
			})
		})
		return ctrl
	}

	appController2 := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test").Version("2")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})
		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController1, appController2},
		})

		return appModule
	}

	return module
}

func Test_VersionURI(t *testing.T) {
	app := core.CreateFactory(AppVersionModule())
	app.SetGlobalPrefix("/api")
	app.EnableVersioning(core.VersionOptions{
		Type: core.URIVersion,
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/v1")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1", res.Data)

	resp2, err := testClient.Get(testServer.URL + "/api/test/v2")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data2, &res2)
	require.Nil(t, err)
	require.Equal(t, "2", res2.Data)
}

func Test_VersionHeader(t *testing.T) {
	app := core.CreateFactory(AppVersionModule())
	app.SetGlobalPrefix("/api")
	app.EnableVersioning(core.VersionOptions{
		Type:   core.HeaderVersion,
		Header: "X-Version",
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)
	req.Header.Set("X-Version", "1")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1", res.Data)

	req.Header.Set("X-Version", "2")
	resp2, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data2, &res2)
	require.Nil(t, err)
	require.Equal(t, "2", res2.Data)

	req.Header.Set("X-Version", "3")
	resp3, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp3.StatusCode)
}

func Test_VersionMedia(t *testing.T) {
	app := core.CreateFactory(AppVersionModule())
	app.SetGlobalPrefix("/api")
	app.EnableVersioning(core.VersionOptions{
		Type: core.MediaTypeVersion,
		Key:  "v=",
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json; v=1")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1", res.Data)

	req.Header.Set("Accept", "application/json; v=2; charset=utf-8")
	resp2, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data2, &res2)
	require.Nil(t, err)
	require.Equal(t, "2", res2.Data)
}

func Test_VersionCustom(t *testing.T) {
	app := core.CreateFactory(AppVersionModule())
	app.SetGlobalPrefix("/api")
	app.EnableVersioning(core.VersionOptions{
		Type: core.CustomVersion,
		Extractor: func(r *http.Request) string {
			return r.URL.Query().Get("version")
		},
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?version=1")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1", res.Data)

	resp2, err := testClient.Get(testServer.URL + "/api/test?version=2")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data2, &res2)
	require.Nil(t, err)
	require.Equal(t, "2", res2.Data)
}
