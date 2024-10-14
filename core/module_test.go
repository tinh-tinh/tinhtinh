package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RequestModule(t *testing.T) {
	tenantModule := func(module *DynamicModule) *DynamicModule {
		tenant := module.New(NewModuleOptions{
			Scope: Request,
		})

		tenant.NewProvider(ProviderOptions{
			Name: "tenant",
			Factory: func(param ...interface{}) interface{} {
				req := param[0].(*http.Request)
				return req.Header.Get("x-tenant")
			},
			Inject: []Provide{REQUEST},
		})
		tenant.Export("tenant")
		return tenant
	}

	userProvider := func(module *DynamicModule) *DynamicProvider {
		provider := module.NewProvider(ProviderOptions{
			Name: "user",
			Factory: func(param ...interface{}) interface{} {
				return fmt.Sprintf("%vUser", param[0])
			},
			Inject: []Provide{Provide("tenant")},
		})
		return provider
	}

	userController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("user")

		ctrl.Get("", func(ctx Ctx) error {
			data := ctrl.Inject(Provide("user"))
			return ctx.JSON(Map{
				"data": data,
			})
		})
		return ctrl
	}

	userModule := func(module *DynamicModule) *DynamicModule {
		user := module.New(NewModuleOptions{
			Scope:       Request,
			Controllers: []Controller{userController},
			Providers:   []Provider{userProvider},
		})
		return user
	}

	appModule := func() *DynamicModule {
		app := NewModule(NewModuleOptions{
			Imports: []Module{tenantModule, userModule},
		})

		return app
	}

	app := CreateFactory(appModule)
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
}

func Test_Controller(t *testing.T) {
	provider := func(module *DynamicModule) *DynamicProvider {
		provider := module.NewProvider(ProviderOptions{
			Name:  "sub",
			Value: "Sub",
		})
		return provider
	}

	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("sub")

		ctrl.Get("", func(ctx Ctx) error {
			data := ctrl.Inject(Provide("sub"))
			return ctx.JSON(Map{
				"data": data,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		module := NewModule(NewModuleOptions{})
		module.Controllers(controller)
		module.Providers(provider)

		return module
	}

	app := CreateFactory(module)
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
