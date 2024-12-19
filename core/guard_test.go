package core_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const Key core.CtxKey = "key"

func TestGuard(t *testing.T) {
	guardInCtrl := func(ctrl core.RefProvider, ctx core.Ctx) bool {
		return ctx.Query("ctrl") == "value"
	}

	guardInModule := func(module core.RefProvider, ctx core.Ctx) bool {
		return ctx.Query("module") == "value"
	}

	guardWithCtx := func(ctrl core.RefProvider, ctx core.Ctx) bool {
		ctx.Set(Key, ctx.Query("ctx"))
		return true
	}

	authCtrl := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Guard(guardInCtrl).Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "1",
			})
		})

		ctrl.Get("module", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})

		ctrl.Get("ctx", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Key),
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{authCtrl},
			Guards:      []core.Guard{guardInModule},
		}).Guard(guardWithCtx)

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/module")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/module?module=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"2"}`, string(data))

	resp, err = testClient.Get(testServer.URL + "/api/test?module=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?module=value&ctrl=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"1"}`, string(data))

	resp, err = testClient.Get(testServer.URL + "/api/test/ctx?module=value&ctx=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"value"}`, string(data))
}
