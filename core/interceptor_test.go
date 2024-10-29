package core

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Transform(ctx *Ctx) CallHandler {
	fmt.Println("Before ...")
	now := time.Now()
	return func(data Map) Map {
		res := make(Map)
		for key, val := range data {
			if val != nil {
				res[key] = val
			}
		}
		fmt.Printf("After ...%vns\n", time.Since(now).Nanoseconds())
		return res
	}
}

func Test_Interceptor(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Interceptor(Transform).Get("", func(ctx Ctx) error {
			return ctx.JSON(Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
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
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"ok","total":10}`, string(data))
}

func Test_ParseInterceptorModule(t *testing.T) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) error {
			return ctx.JSON(Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
			})
		})

		return ctrl
	}

	appModule := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController},
			Interceptor: Transform,
		})

		return appModule
	}

	app := CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"ok","total":10}`, string(data))
}
