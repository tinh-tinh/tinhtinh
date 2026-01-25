package logger_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

func TestMiddleware(t *testing.T) {
	appController := func(module core.Module) core.Controller {
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

	appModule := func() core.Module {
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

	time.Sleep(1 * time.Second)
}

func TestFormatTemplates(t *testing.T) {
	t.Run("Dev format", func(t *testing.T) {
		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			Format: logger.Dev,
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		time.Sleep(500 * time.Millisecond)
	})

	t.Run("Common format", func(t *testing.T) {
		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
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

		time.Sleep(500 * time.Millisecond)
	})

	t.Run("Combined format", func(t *testing.T) {
		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			Format: logger.Combined,
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		req, _ := http.NewRequest("GET", testServer.URL+"/api/test", nil)
		req.Header.Set("User-Agent", "TestAgent/1.0")
		req.Header.Set("Referer", "http://example.com")

		resp, err := testClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		time.Sleep(500 * time.Millisecond)
	})

	t.Run("default format (empty)", func(t *testing.T) {
		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			// Format not specified - should default to Dev
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		time.Sleep(500 * time.Millisecond)
	})
}

func TestCustomFormatter(t *testing.T) {
	t.Run("basic custom formatter", func(t *testing.T) {
		var capturedLog string

		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			CustomFormatter: func(ctx logger.LogContext) string {
				capturedLog = fmt.Sprintf("[CUSTOM] %s %s -> %d (%s)",
					ctx.Request.Method,
					ctx.Request.URL.Path,
					ctx.StatusCode,
					ctx.Latency,
				)
				return capturedLog
			},
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		time.Sleep(500 * time.Millisecond)

		// Verify custom formatter was called
		require.Contains(t, capturedLog, "[CUSTOM]")
		require.Contains(t, capturedLog, "GET")
		require.Contains(t, capturedLog, "/api/test")
		require.Contains(t, capturedLog, "200")
	})

	t.Run("custom formatter with response headers", func(t *testing.T) {
		var capturedHeaders http.Header

		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				ctx.Res().Header().Set("X-Custom-Header", "test-value")
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			CustomFormatter: func(ctx logger.LogContext) string {
				capturedHeaders = ctx.ResponseHeaders
				return fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path)
			},
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		time.Sleep(500 * time.Millisecond)

		// Verify response headers were captured
		require.NotNil(t, capturedHeaders)
	})

	t.Run("custom formatter with different status codes", func(t *testing.T) {
		var lastStatusCode int

		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("ok", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})

			ctrl.Get("notfound", func(ctx core.Ctx) error {
				return ctx.Status(http.StatusNotFound).JSON(core.Map{"error": "not found"})
			})

			ctrl.Get("error", func(ctx core.Ctx) error {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "internal error"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			CustomFormatter: func(ctx logger.LogContext) string {
				lastStatusCode = ctx.StatusCode
				return fmt.Sprintf("Status: %d", ctx.StatusCode)
			},
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		// Test 200
		resp, err := testClient.Get(testServer.URL + "/api/test/ok")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		time.Sleep(200 * time.Millisecond)
		require.Equal(t, 200, lastStatusCode)

		// Test 404
		resp, err = testClient.Get(testServer.URL + "/api/test/notfound")
		require.Nil(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		time.Sleep(200 * time.Millisecond)
		require.Equal(t, 404, lastStatusCode)

		// Test 500
		resp, err = testClient.Get(testServer.URL + "/api/test/error")
		require.Nil(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		time.Sleep(200 * time.Millisecond)
		require.Equal(t, 500, lastStatusCode)
	})

	t.Run("custom formatter with latency check", func(t *testing.T) {
		var capturedLatency time.Duration

		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")

			ctrl.Get("", func(ctx core.Ctx) error {
				time.Sleep(50 * time.Millisecond) // Simulate some processing
				return ctx.JSON(core.Map{"data": "ok"})
			})

			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			CustomFormatter: func(ctx logger.LogContext) string {
				capturedLatency = ctx.Latency
				return fmt.Sprintf("Latency: %s", ctx.Latency)
			},
			Rotate: true,
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()

		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		time.Sleep(500 * time.Millisecond)

		// Verify latency was captured and is reasonable (>= 50ms due to sleep)
		require.GreaterOrEqual(t, capturedLatency.Milliseconds(), int64(50))
	})
}

func TestSkipPaths(t *testing.T) {
	t.Run("skip paths", func(t *testing.T) {
		var logCalled bool

		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")
			ctrl.Get("health", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"status": "ok"})
			})
			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			CustomFormatter: func(ctx logger.LogContext) string {
				logCalled = true
				return "log"
			},
			SkipPaths: []string{"/api/test/health"},
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()
		resp, err := testClient.Get(testServer.URL + "/api/test/health")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.False(t, logCalled, "Logger should not be called for skipped path")
	})

	t.Run("non-skipped paths", func(t *testing.T) {
		var logCalled bool

		appController := func(module core.Module) core.Controller {
			ctrl := module.NewController("test")
			ctrl.Get("normal", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{"data": "ok"})
			})
			return ctrl
		}

		appModule := func() core.Module {
			return core.NewModule(core.NewModuleOptions{
				Controllers: []core.Controllers{appController},
			})
		}

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")
		app.Use(logger.Handler(logger.MiddlewareOptions{
			CustomFormatter: func(ctx logger.LogContext) string {
				logCalled = true
				return "log"
			},
			SkipPaths: []string{"/api/test/health"},
		}))

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()

		testClient := testServer.Client()
		resp, err := testClient.Get(testServer.URL + "/api/test/normal")
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.True(t, logCalled, "Logger should be called for non-skipped path")
	})
}
