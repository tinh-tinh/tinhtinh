package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Route(t *testing.T) {
	route := ParseRoute("GET /test")

	route.SetPrefix("api")
	require.Equal(t, "GET /api/test", route.GetPath())

	route = ParseRoute("GET /user/{id}")
	require.Equal(t, "GET /user/{id}", route.GetPath())

	route = ParseRoute("GET /user/{id}/edit")
	route.SetPrefix("admin")
	require.Equal(t, "GET /admin/user/{id}/edit", route.GetPath())
}

func Test_IfSlashPrefixString(t *testing.T) {
	require.Equal(t, "", IfSlashPrefixString(""))
	require.Equal(t, "/", IfSlashPrefixString("/"))
	require.Equal(t, "/api", IfSlashPrefixString("api"))
	require.Equal(t, "/api", IfSlashPrefixString("/api"))
	require.Equal(t, "/api", IfSlashPrefixString("API"))
	require.Equal(t, "/api", IfSlashPrefixString("a pi"))
}

func Test_registerRoutes(t *testing.T) {
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

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err := http.NewRequest("PATCH", testServer.URL+"/api/test/1", nil)
	require.Nil(t, err)
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("PUT", testServer.URL+"/api/test/1", nil)
	require.Nil(t, err)
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("DELETE", testServer.URL+"/api/test/1", nil)
	require.Nil(t, err)
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
