package core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseGuardCtrl(t *testing.T) {
	guard := func(ctrl *DynamicController, ctx *Ctx) bool {
		return ctx.Query("key") == "value"
	}

	authCtrl := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Guard(guard).Get("", func(ctx Ctx) error {
			return ctx.JSON(Map{
				"data": "1",
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{authCtrl},
		})

		return appModule
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?key=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?key=abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func Test_ParseGuardModule(t *testing.T) {
	guard := func(module *DynamicModule, ctx Ctx) bool {
		return ctx.Query("key") == "value"
	}

	authCtrl := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) error {
			return ctx.JSON(Map{
				"data": "1",
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{authCtrl},
		})

		appModule.Guard(guard)

		return appModule
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?key=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?key=abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

const Key CtxKey = "key"

func Test_Ctx_Guard(t *testing.T) {
	guard := func(ctrl *DynamicController, ctx *Ctx) bool {
		ctx.Set(Key, "value")
		return true
	}

	authCtrl := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Guard(guard).Get("", func(ctx Ctx) error {
			return ctx.JSON(Map{
				"data": ctx.Get(Key),
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{authCtrl},
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

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "value", res.Data)
}
