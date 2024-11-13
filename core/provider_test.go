package core_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func ChildProvider(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name:  "child",
		Value: "child",
	})
	return provider
}

func ChildModule(module *core.DynamicModule) *core.DynamicModule {
	childModule := module.New(core.NewModuleOptions{
		Scope:     core.Global,
		Providers: []core.Provider{ChildProvider},
		Exports:   []core.Provider{ChildProvider},
	})

	return childModule
}

func AppController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("test")
	ctrl.Get("/", func(ctx core.Ctx) error {
		name := ctrl.Inject("child")
		return ctx.JSON(core.Map{
			"data": name,
		})
	})
	return ctrl
}

func AppModule() *core.DynamicModule {
	module := core.NewModule(core.NewModuleOptions{
		Scope:       core.Global,
		Imports:     []core.Module{ChildModule},
		Controllers: []core.Controller{AppController},
	})

	return module
}

func Test_NewProvider(t *testing.T) {
	app := core.CreateFactory(AppModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "child", res.Data)
}

// func Test_getExports(t *testing.T) {
// 	childPrivateProvider := func(module *core.DynamicModule) *core.DynamicProvider {
// 		provider := module.NewProvider(core.ProviderOptions{
// 			Name:  "private",
// 			Value: "private",
// 		})
// 		return provider
// 	}

// 	childPublicProvider := func(module *core.DynamicModule) *core.DynamicProvider {
// 		provider := module.NewProvider(core.ProviderOptions{
// 			Name:  "public",
// 			Value: "public",
// 		})
// 		return provider
// 	}

// 	childModule := func(module *core.DynamicModule) *core.DynamicModule {
// 		childModule := module.New(core.NewModuleOptions{
// 			Providers: []core.Provider{childPrivateProvider, childPublicProvider},
// 			Exports:   []core.Provider{childPublicProvider},
// 		})

// 		return childModule
// 	}

// 	module := core.NewModule(core.NewModuleOptions{
// 		Imports: []core.Module{childModule},
// 	})
// 	providers := module.getExports()
// 	require.Equal(t, 1, len(providers))
// 	require.Equal(t, core.Provide("public"), providers[0].Name)
// }

// func Test_getRequest(t *testing.T) {
// 	reqModule := func(module *core.DynamicModule) *core.DynamicModule {
// 		req := module.New(core.NewModuleOptions{})
// 		req.NewProvider(core.ProviderOptions{
// 			Scope: core.Request,
// 			Name:  "req",
// 			Factory: func(param ...interface{}) interface{} {
// 				return param[0]
// 			},
// 			Inject: []core.Provide{core.REQUEST},
// 		})
// 		req.Export("req")
// 		req.NewProvider(core.ProviderOptions{
// 			Name:  "req2",
// 			Value: 3,
// 		})
// 		req.Export("req2")

// 		fmt.Println(req.getExports())
// 		return req
// 	}

// 	globalModule := func(module *core.DynamicModule) *core.DynamicModule {
// 		global := module.New(core.NewModuleOptions{
// 			Scope: core.Global,
// 		})
// 		global.NewProvider(core.ProviderOptions{
// 			Name:  "global",
// 			Value: "global",
// 		})

// 		global.Export("global")
// 		return global
// 	}

// 	module := core.NewModule(core.NewModuleOptions{
// 		Imports: []core.Module{reqModule, globalModule},
// 	})
// 	providers := module.getRequest()
// 	fmt.Println(providers)
// 	require.Equal(t, 1, len(providers))
// 	require.Equal(t, core.Provide("req"), providers[0].Name)
// }

// func Test_appendProvider(t *testing.T) {
// 	module := core.NewModule(core.NewModuleOptions{
// 		Scope: core.Global,
// 	})

// 	provider := module.NewProvider(core.ProviderOptions{
// 		Name:  "test",
// 		Value: "test",
// 	})

// 	module.appendProvider(provider)
// 	require.Equal(t, 1, len(module.DataProviders))
// }

func Test_FactoryProvider(t *testing.T) {
	rootModule := func(module *core.DynamicModule) *core.DynamicModule {
		root := module.New(core.NewModuleOptions{
			Scope: core.Global,
		})
		root.NewProvider(core.ProviderOptions{
			Name:  "root",
			Value: "root",
		})
		root.Export("root")

		return root
	}

	childModule := func(module *core.DynamicModule) *core.DynamicModule {
		child := module.New(core.NewModuleOptions{
			Scope: core.Global,
		})
		child.NewProvider(core.ProviderOptions{
			Name: "child",
			Factory: func(param ...interface{}) interface{} {
				return fmt.Sprintf("%vChild", param[0])
			},
			Inject: []core.Provide{core.Provide("root")},
		})
		child.Export("child")

		return child
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{rootModule, childModule},
	})

	require.Equal(t, "rootChild", module.Ref("child"))
}

func Test_RequestProvider(t *testing.T) {
	const (
		TENANT  core.Provide = "TENANT"
		SERVICE core.Provide = "SERVICE"
	)
	type RequestProvider struct {
		Name string
	}
	service := func(module *core.DynamicModule) *core.DynamicProvider {
		prd := module.NewProvider(core.ProviderOptions{
			Scope: core.Request,
			Name:  SERVICE,
			Factory: func(param ...interface{}) interface{} {
				return &RequestProvider{
					Name: "model" + param[0].(string),
				}
			},
			Inject: []core.Provide{TENANT},
		})

		return prd
	}

	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("/", func(ctx core.Ctx) error {
			service := module.Ref(SERVICE).(*RequestProvider)
			return ctx.JSON(core.Map{
				"data": service.Name,
			})
		})

		return ctrl
	}

	tenantModule := func(module *core.DynamicModule) *core.DynamicModule {
		tenant := module.New(core.NewModuleOptions{
			Scope: core.Global,
		})

		tenant.NewProvider(core.ProviderOptions{
			Name: TENANT,
			Factory: func(param ...interface{}) interface{} {
				return param[0].(*http.Request).Header.Get("x-tenant")
			},
			Inject: []core.Provide{core.REQUEST},
		})
		tenant.Export(TENANT)

		return tenant
	}

	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{
			Imports:     []core.Module{tenantModule},
			Controllers: []core.Controller{controller},
			Providers:   []core.Provider{service},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test/", nil)
	require.Nil(t, err)
	req.Header.Set("x-tenant", "test")
	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "modeltest", res.Data)

	req.Header.Set("x-tenant", "test2")
	resp2, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data2, &res2)
	require.Nil(t, err)
	require.Equal(t, "modeltest2", res2.Data)
}
