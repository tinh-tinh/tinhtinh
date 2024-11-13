package core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CustomCtx(t *testing.T) {
	tenant := CreateWrapper(func(data bool, ctx Ctx) string {
		if data {
			return "master"
		}
		return ctx.Req().Header.Get("x-tenant-id")
	})
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", tenant.Handler(true, func(wCtx WrappedCtx[string]) error {
			return wCtx.JSON(Map{
				"data": wCtx.Data,
			})
		}))

		ctrl.Get("tenant", tenant.Handler(false, func(wCtx WrappedCtx[string]) error {
			return wCtx.JSON(Map{
				"data": wCtx.Data,
			})
		}))

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController},
		})

		return appModule
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "babadook")
	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "master", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/test/tenant", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "babadook")
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data, &res2)
	require.Nil(t, err)
	require.Equal(t, "babadook", res2.Data)
}

func Test_Middleware_CustomCtx(t *testing.T) {
	tenant := CreateWrapper(func(data bool, ctx Ctx) string {
		if data {
			return "master"
		}
		return ctx.Req().Header.Get("x-tenant-id")
	})

	guard := func(ctrl *DynamicController, ctx *Ctx) bool {
		return ctx.Query("key") == "value"
	}

	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Guard(guard).Get("", tenant.Handler(true, func(wCtx WrappedCtx[string]) error {
			return wCtx.JSON(Map{
				"data": wCtx.Data,
			})
		}))

		ctrl.Get("tenant", tenant.Handler(false, func(wCtx WrappedCtx[string]) error {
			return wCtx.JSON(Map{
				"data": wCtx.Data,
			})
		}))

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController},
		})

		return appModule
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?key=haha")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	req, err := http.NewRequest("GET", testServer.URL+"/api/test?key=value", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "babadook")
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "master", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/test/tenant", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "babadook")
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data, &res2)
	require.Nil(t, err)
	require.Equal(t, "babadook", res2.Data)
}
