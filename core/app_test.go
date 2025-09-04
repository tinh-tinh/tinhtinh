package core_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Exception(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("bad-request", func(ctx core.Ctx) error {
			return common.BadRequestException(ctx.Res(), "bad request")
		})

		ctrl.Get("unauthorized", func(ctx core.Ctx) error {
			return common.UnauthorizedException(ctx.Res(), "unauthorized")
		})

		ctrl.Get("forbidden", func(ctx core.Ctx) error {
			return common.ForbiddenException(ctx.Res(), "forbidden")
		})

		ctrl.Get("not-found", func(ctx core.Ctx) error {
			return common.NotFoundException(ctx.Res(), "not found")
		})

		ctrl.Get("method-not-allowed", func(ctx core.Ctx) error {
			return common.NotAllowedException(ctx.Res(), "method not allowed")
		})

		ctrl.Get("conflict", func(ctx core.Ctx) error {
			return common.ConflictException(ctx.Res(), "conflict")
		})

		ctrl.Get("internal-server-error", func(ctx core.Ctx) error {
			return common.InternalServerException(ctx.Res(), "internal server error")
		})

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
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "data",
			})
		})

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

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resp, err := testClient.Get(testServer.URL + "/api/test")
		require.Nil(b, err)
		require.Equal(b, http.StatusOK, resp.StatusCode)
	}
}

func Test_Timeout(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			time.Sleep(3 * time.Second)
			return ctx.JSON(core.Map{
				"data": "data",
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module, core.AppOptions{
		Timeout: 1 * time.Second,
	})
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func Test_Listen(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

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

	require.NotPanics(t, func() {
		go func() {
			app.Listen(3000)
		}()
	})
}

func Test_CustomValidate(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		type QueryId struct {
			ID string `query:"id"`
		}
		ctrl.Pipe(core.QueryParser[QueryId]{}).Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module, core.AppOptions{
		CustomValidation: func(val any) error {
			return errors.New("Alway error")
		},
	})
	app.SetGlobalPrefix("api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?id=2")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_Encode(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		type BodyInput struct {
			ID string `json:"id"`
		}
		ctrl.Pipe(core.BodyParser[BodyInput]{}).Post("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Body(),
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module, core.AppOptions{
		Encoder: json.Marshal,
		Decoder: json.Unmarshal,
	})
	app.SetGlobalPrefix("api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"id": "1"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)

	require.NotNil(t, res.Data)
}
