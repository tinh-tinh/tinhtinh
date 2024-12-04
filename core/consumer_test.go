package core_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Consumer(t *testing.T) {
	const (
		Tenant   core.CtxKey = "tenant"
		Location core.CtxKey = "location"
	)

	tenantMiddleware := func(ctx core.Ctx) error {
		tenant := ctx.Headers("x-tenant-id")
		if tenant != "" {
			ctx.Set(Tenant, tenant)
		}
		return ctx.Next()
	}

	locationMiddleware := func(ctx core.Ctx) error {
		location := ctx.Headers("x-location-id")
		if location != "" {
			ctx.Set(Location, location)
		}
		return ctx.Next()
	}

	userMiddleware := func(ctx core.Ctx) error {
		user := ctx.Headers("x-user-id")
		if user != "" {
			ctx.Set(Tenant, user)
		}
		return ctx.Next()
	}

	userController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("user")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Tenant),
			})
		})

		ctrl.Get("location", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Location),
			})
		})

		ctrl.Get("special", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "special",
			})
		})

		return ctrl
	}

	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		user := module.New(core.NewModuleOptions{
			Controllers: []core.Controller{userController},
		})

		user.Consumer(core.NewConsumer().Apply(userMiddleware).Include(core.RoutesPath{
			Path: "/user", Method: http.MethodGet,
		}, core.RoutesPath{
			Path: "*", Method: http.MethodGet,
		}, core.RoutesPath{
			Path: "/user/location", Method: core.MethodAll,
		}, core.RoutesPath{
			Path: "/user/special", Method: http.MethodGet,
		}))

		return user
	}

	postController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("post")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Tenant),
			})
		})

		ctrl.Get("exclude", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Location),
			})
		})

		ctrl.Post("special", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "special",
			})
		})

		return ctrl
	}

	postModule := func(module *core.DynamicModule) *core.DynamicModule {
		post := module.New(core.NewModuleOptions{
			Controllers: []core.Controller{postController},
		})

		post.Consumer(core.NewConsumer().Apply(userMiddleware).Exclude(core.RoutesPath{
			Path: "*", Method: http.MethodGet,
		}, core.RoutesPath{
			Path: "/post/exclude", Method: core.MethodAll,
		}, core.RoutesPath{
			Path: "/post/special", Method: http.MethodPost,
		}))

		return post
	}

	appModule := func() *core.DynamicModule {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{userModule, postModule},
		})

		app.Consumer(core.NewConsumer().Apply(tenantMiddleware).Include(core.RoutesPath{
			Path: "*", Method: core.MethodAll,
		}))

		app.Consumer(core.NewConsumer().Apply(locationMiddleware).Exclude(core.RoutesPath{
			Path: "/post/exclude", Method: core.MethodAll,
		}))

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/user", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "test")
	req.Header.Set("x-location-id", "test2")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/user/location", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "test")
	req.Header.Set("x-location-id", "test2")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test2", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/post", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "test")
	req.Header.Set("x-location-id", "test2")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/post/exclude", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "test")
	req.Header.Set("x-location-id", "test2")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Nil(t, res.Data)
}
