package core_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_DefaultErrorHandler(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			panic(exception.ThrowHttp("test", http.StatusInternalServerError))
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	type ErrorData struct {
		StatusCode int    `json:"statusCode"`
		Error      string `json:"error"`
		Timestamp  string `json:"timestamp"`
		Path       string `json:"path"`
	}
	var errData ErrorData
	err = json.Unmarshal(data, &errData)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, errData.StatusCode)
	require.Equal(t, "test", errData.Error)
	require.Equal(t, "/api/test", errData.Path)
}

func Test_ErrorHandler(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return errors.New("Error")
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
