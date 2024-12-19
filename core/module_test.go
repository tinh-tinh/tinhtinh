package core_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_RequestModule(t *testing.T) {
	tenantModule := func(module core.Module) core.Module {
		tenant := module.New(core.NewModuleOptions{
			Scope: core.Request,
		})

		tenant.NewProvider(core.ProviderOptions{
			Name: "tenant",
			Factory: func(param ...interface{}) interface{} {
				req := param[0].(*http.Request)
				return req.Header.Get("x-tenant")
			},
			Inject: []core.Provide{core.REQUEST},
		})
		tenant.Export("tenant")
		return tenant
	}

	userProvider := func(module core.Module) core.Provider {
		provider := module.NewProvider(core.ProviderOptions{
			Name: "user",
			Factory: func(param ...interface{}) interface{} {
				return fmt.Sprintf("%vUser", param[0])
			},
			Inject: []core.Provide{core.Provide("tenant")},
		})
		return provider
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("user")

		ctrl.Get("", func(ctx core.Ctx) error {
			fmt.Println("in here")
			data := ctrl.Ref(core.Provide("user"), ctx)
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("panic", func(ctx core.Ctx) error {
			data := ctrl.Ref(core.Provide("user"))
			return ctx.JSON(core.Map{
				"data": data,
			})
		})
		return ctrl
	}

	userModule := func(module core.Module) core.Module {
		user := module.New(core.NewModuleOptions{
			Scope:       core.Request,
			Controllers: []core.Controllers{userController},
			Providers:   []core.Providers{userProvider},
		})
		return user
	}

	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{tenantModule, userModule},
		})

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	req, err := http.NewRequest("GET", testServer.URL+"/api/user", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant", "1")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1User", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/user", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant", "2")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "2User", res.Data)

	// req, err = http.NewRequest("GET", testServer.URL+"/api/user/panic", nil)
	// require.Nil(t, err)
	// req.Header.Set("x-tenant", "3")

	// resp, err = testClient.Do(req)
	// require.Nil(t, err)
	// require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Controller(t *testing.T) {
	provider := func(module core.Module) core.Provider {
		provider := module.NewProvider(core.ProviderOptions{
			Name:  "sub",
			Value: "Sub",
		})
		return provider
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("sub")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctrl.Ref(core.Provide("sub"))
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		return ctrl
	}

	module := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{})
		module.Controllers(controller)
		module.Providers(provider)

		return module
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	resp, err := testClient.Get(testServer.URL + "/api/sub")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "Sub", res.Data)
}

func Test_Nil(t *testing.T) {
	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{nil},
			Providers:   []core.Providers{nil},
			Imports:     []core.Modules{nil},
			Exports:     []core.Providers{nil},
			Guards:      []core.Guard{nil},
			Middlewares: []core.Middleware{nil},
		})
		return module
	}

	app := appModule()
	prd := app.Ref("abc")
	require.Nil(t, prd)

	require.NotPanics(t, func() {
		_ = core.CreateFactory(appModule)
	})
}
func Test_Import(t *testing.T) {
	const SUB_PROVIDER core.Provide = "sub"
	subModule := func(module core.Module) core.Module {
		sub := module.New(core.NewModuleOptions{})
		sub.NewProvider(core.ProviderOptions{
			Name:  SUB_PROVIDER,
			Value: "haha",
		})
		sub.Export(SUB_PROVIDER)

		return sub
	}

	const PARENT_SERVICE core.Provide = "parent"
	parentService := func(module core.Module) core.Provider {
		s := module.NewProvider(core.ProviderOptions{
			Name: PARENT_SERVICE,
			Factory: func(param ...interface{}) interface{} {
				sub := param[0].(string)
				return sub + "hihi"
			},
			Inject: []core.Provide{SUB_PROVIDER},
		})
		return s
	}

	parentModule := func(module core.Module) core.Module {
		parent := module.New(core.NewModuleOptions{
			Imports:   []core.Modules{subModule},
			Providers: []core.Providers{parentService},
			Exports:   []core.Providers{parentService},
		})

		return parent
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{parentModule},
		})

		return module
	}

	require.NotPanics(t, func() {
		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("api")
	})
}

func Test_LifecycleModule(t *testing.T) {
	const Tenant core.CtxKey = "tenant"
	tenantMiddleware := func(ctx core.Ctx) error {
		tenant := ctx.Headers("x-tenant-id")
		if tenant != "" {
			ctx.Set(Tenant, tenant)
		}
		return ctx.Next()
	}

	tenantGuard := func(module core.RefProvider, ctx core.Ctx) bool {
		return ctx.Get(Tenant) != nil
	}

	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Tenant),
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Middlewares: []core.Middleware{tenantMiddleware},
			Guards:      []core.Guard{tenantGuard},
			Controllers: []core.Controllers{appController},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)
	req.Header.Set("X-Tenant-Id", "1")
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1", res.Data)
}

func Test_PassMiddlewareModule(t *testing.T) {
	const Tenant core.CtxKey = "tenant"
	tenantMiddleware := func(ctx core.Ctx) error {
		tenant := ctx.Headers("x-tenant-id")
		if tenant != "" {
			ctx.Set(Tenant, tenant)
		}
		return ctx.Next()
	}

	tenantGuard := func(module core.RefProvider, ctx core.Ctx) bool {
		return ctx.Get(Tenant) != nil
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("user")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Get(Tenant),
			})
		})

		return ctrl
	}

	userModule := func(module core.Module) core.Module {
		user := module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{userController},
		})

		return user
	}

	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Imports:     []core.Modules{userModule},
			Middlewares: []core.Middleware{tenantMiddleware},
			Guards:      []core.Guard{tenantGuard},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/user")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	req, err := http.NewRequest("GET", testServer.URL+"/api/user", nil)
	require.Nil(t, err)
	req.Header.Set("X-Tenant-Id", "1")
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "1", res.Data)
}
