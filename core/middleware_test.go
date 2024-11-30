package core_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_CtxContext(t *testing.T) {
	const key core.CtxKey = "key"

	middleware := func(ctx core.Ctx) error {
		ctx.Set(key, "value")
		return ctx.Next()
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(middleware).Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(key),
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "value", res.Data)
}

func Test_Middleware(t *testing.T) {
	const key core.CtxKey = "key"

	middleware := func(ctx core.Ctx) error {
		ctx.Set(key, "value")
		return ctx.Next()
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(key),
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		}).Use(middleware)

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "value", res.Data)
}

func Test_ExceptionMiddleware(t *testing.T) {
	middleware := func(ctx core.Ctx) error {
		panic("error")
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		}).Use(middleware)

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	type ErrorResponse struct {
		Error interface{} `json:"error"`
		Path  string      `json:"path"`
	}
	var res ErrorResponse
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "error", res.Error)
}
