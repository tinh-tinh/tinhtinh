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

func Test_CustomCtx(t *testing.T) {
	tenant := core.CreateWrapper(func(data bool, ctx core.Ctx) string {
		if data {
			return "master"
		}
		return ctx.Req().Header.Get("x-tenant-id")
	})
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", tenant.Handler(true, func(wCtx core.WrappedCtx[string]) error {
			return wCtx.JSON(core.Map{
				"data": wCtx.Data,
			})
		}))

		ctrl.Get("tenant", tenant.Handler(false, func(wCtx core.WrappedCtx[string]) error {
			return wCtx.JSON(core.Map{
				"data": wCtx.Data,
			})
		}))

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
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
	tenant := core.CreateWrapper(func(data bool, ctx core.Ctx) string {
		if data {
			return "master"
		}
		return ctx.Req().Header.Get("x-tenant-id")
	})

	guard := func(ctx core.Ctx) bool {
		return ctx.Query("key") == "value"
	}

	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Guard(guard).Get("", tenant.Handler(true, func(wCtx core.WrappedCtx[string]) error {
			return wCtx.JSON(core.Map{
				"data": wCtx.Data,
			})
		}))

		ctrl.Get("tenant", tenant.Handler(false, func(wCtx core.WrappedCtx[string]) error {
			return wCtx.JSON(core.Map{
				"data": wCtx.Data,
			})
		}))

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
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
