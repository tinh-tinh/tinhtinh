package core_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Route(t *testing.T) {
	route := core.ParseRoute("GET /test")

	route.SetPrefix("api")
	require.Equal(t, "GET /api/test", route.GetPath())

	route = core.ParseRoute("GET /user/{id}")
	require.Equal(t, "GET /user/{id}", route.GetPath())

	route = core.ParseRoute("GET /user/{id}/edit")
	route.SetPrefix("admin")
	require.Equal(t, "GET /admin/user/{id}/edit", route.GetPath())

	route = core.ParseRoute("GET /")
	require.Equal(t, "GET /", route.GetPath())
}

func Test_IfSlashPrefixString(t *testing.T) {
	require.Equal(t, "", core.IfSlashPrefixString(""))
	require.Equal(t, "/", core.IfSlashPrefixString("/"))
	require.Equal(t, "/api", core.IfSlashPrefixString("api"))
	require.Equal(t, "/api", core.IfSlashPrefixString("/api"))
	require.Equal(t, "/api", core.IfSlashPrefixString("API"))
	require.Equal(t, "/api", core.IfSlashPrefixString("a pi"))
}

func Test_registerRoutes(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "1",
			})
		})

		ctrl.Post("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})

		ctrl.Patch("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "3",
			})
		})

		ctrl.Put("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "4",
			})
		})

		ctrl.Delete("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "5",
			})
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
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"1"}`, string(data))

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"2"}`, string(data))

	req, err := http.NewRequest("PATCH", testServer.URL+"/api/test/1", nil)
	require.Nil(t, err)
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"3"}`, string(data))

	req, err = http.NewRequest("PUT", testServer.URL+"/api/test/1", nil)
	require.Nil(t, err)
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"4"}`, string(data))

	req, err = http.NewRequest("DELETE", testServer.URL+"/api/test/1", nil)
	require.Nil(t, err)
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"5"}`, string(data))
}
