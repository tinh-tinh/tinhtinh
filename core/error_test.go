package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DefaultErrorHandler(t *testing.T) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) error {
			panic(errors.New("Error"))
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	fmt.Println(string(data))
}

func Test_ErrorHandler(t *testing.T) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) error {
			return errors.New("Error")
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController},
		})

		return appModule
	}

	app := CreateFactory(module, AppOptions{
		ErrorHandler: func(err error, ctx Ctx) error {
			return ctx.Status(http.StatusInternalServerError).JSON(Map{
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
