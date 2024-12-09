package logger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

func TestMiddleware(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("success", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("redirect", func(ctx core.Ctx) error {
			return ctx.Status(http.StatusMovedPermanently).JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("fail", func(ctx core.Ctx) error {
			return ctx.Status(http.StatusNotFound).JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("error", func(ctx core.Ctx) error {
			return ctx.Status(http.StatusInternalServerError).JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")
	app.Use(logger.Handler(logger.MiddlewareOptions{
		SeparateBaseStatus: true,
		Format:             "${method} ${path} ${status} ${latency}",
		Rotate:             true,
	}))

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/success")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/redirect")
	require.Nil(t, err)
	require.Equal(t, http.StatusMovedPermanently, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/fail")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/error")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDev(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")
	app.Use(logger.Handler(logger.MiddlewareOptions{
		Format: logger.Common,
		Rotate: true,
	}))

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
