package core_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_DefaultErrorHandler(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			panic(errors.New("Error"))
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
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
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	fmt.Println(string(data))
}

func Test_ErrorHandler(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return errors.New("Error")
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module, core.AppOptions{
		ErrorHandler: func(err error, ctx core.Ctx) error {
			return ctx.Status(http.StatusInternalServerError).JSON(core.Map{
				"message": err.Error(),
			})
		},
	})
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	require.Equal(t, "{\"message\":\"Error\"}", string(data))
}
