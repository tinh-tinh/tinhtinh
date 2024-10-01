package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/middleware/cors"
)

func Test_EnableCors(t *testing.T) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "1",
			})
		})

		ctrl.Post("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "2",
			})
		})

		ctrl.Patch("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "3",
			})
		})

		ctrl.Put("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "4",
			})
		})

		ctrl.Delete("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "5",
			})
		})
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
	app.EnableCors(cors.Options{
		AllowedMethods: []string{"POST"},
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func Test_Exception(t *testing.T) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("bad-request", func(ctx Ctx) {
			common.BadRequestException(ctx.Res(), "bad request")
		})

		ctrl.Get("unauthorized", func(ctx Ctx) {
			common.UnauthorizedException(ctx.Res(), "unauthorized")
		})

		ctrl.Get("forbidden", func(ctx Ctx) {
			common.ForbiddenException(ctx.Res(), "forbidden")
		})

		ctrl.Get("not-found", func(ctx Ctx) {
			common.NotFoundException(ctx.Res(), "not found")
		})

		ctrl.Get("method-not-allowed", func(ctx Ctx) {
			common.NotAllowedException(ctx.Res(), "method not allowed")
		})

		ctrl.Get("conflict", func(ctx Ctx) {
			common.ConflictException(ctx.Res(), "conflict")
		})

		ctrl.Get("internal-server-error", func(ctx Ctx) {
			common.InternalServerException(ctx.Res(), "internal server error")
		})

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

	resp, err := testClient.Get(testServer.URL + "/api")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/bad-request")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/unauthorized")
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/forbidden")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/not-found")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/method-not-allowed")
	require.Nil(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/conflict")
	require.Nil(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/internal-server-error")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Benchmark_App(b *testing.B) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			data := make(map[string]string)

			for i := 0; i < b.N; i++ {
				data[fmt.Sprintf("%d", i)] = fmt.Sprintf("%d", i)
			}

			ctx.JSON(Map{
				"data": data,
			})

			runtime.GC()
		})

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

	for i := 0; i < b.N; i++ {
		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(b, err)
		require.Equal(b, http.StatusOK, resp.StatusCode)
	}
}
