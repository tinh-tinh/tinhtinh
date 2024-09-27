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

func ChildProvider(module *DynamicModule) *DynamicProvider {
	provider := module.NewProvider(ProviderOptions{
		Name:  "child",
		Value: "child",
	})
	return provider
}

func ChildModule(module *DynamicModule) *DynamicModule {
	childModule := module.New(NewModuleOptions{
		Scope:     Global,
		Providers: []Provider{ChildProvider},
		Exports:   []Provider{ChildProvider},
	})

	return childModule
}

func AppController(module *DynamicModule) *DynamicController {
	ctrl := module.NewController("test")
	ctrl.Get("/", func(ctx Ctx) {
		name := ctrl.Inject("child")
		ctx.JSON(Map{
			"data": name,
		})
	})
	return ctrl
}

func AppModule() *DynamicModule {
	module := NewModule(NewModuleOptions{
		Scope:       Global,
		Imports:     []Module{ChildModule},
		Controllers: []Controller{AppController},
	})

	return module
}

func Test_NewProvider(t *testing.T) {
	app := CreateFactory(AppModule, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "child", res.Data)
}

func Test_getExports(t *testing.T) {
	childPrivateProvider := func(module *DynamicModule) *DynamicProvider {
		provider := module.NewProvider(ProviderOptions{
			Name:  "private",
			Value: "private",
		})
		return provider
	}

	childPublicProvider := func(module *DynamicModule) *DynamicProvider {
		provider := module.NewProvider(ProviderOptions{
			Name:  "public",
			Value: "public",
		})
		return provider
	}

	childModule := func(module *DynamicModule) *DynamicModule {
		childModule := module.New(NewModuleOptions{
			Providers: []Provider{childPrivateProvider, childPublicProvider},
			Exports:   []Provider{childPublicProvider},
		})

		return childModule
	}

	module := NewModule(NewModuleOptions{
		Imports: []Module{childModule},
	})
	providers := module.getExports()
	require.Equal(t, 1, len(providers))
	require.Equal(t, Provide("public"), providers[0].Name)
}

func Test_getRequest(t *testing.T) {
	reqModule := func(module *DynamicModule) *DynamicModule {
		req := module.New(NewModuleOptions{
			Scope: Request,
		})
		req.NewProvider(ProviderOptions{
			Name: "req",
			Factory: func(param ...interface{}) interface{} {
				return param[0]
			},
			Inject: []Provide{REQUEST},
		})
		req.Export("req")
		return req
	}

	globalModule := func(module *DynamicModule) *DynamicModule {
		global := module.New(NewModuleOptions{
			Scope: Global,
		})
		global.NewProvider(ProviderOptions{
			Name:  "global",
			Value: "global",
		})

		global.Export("global")
		return global
	}

	module := NewModule(NewModuleOptions{
		Imports: []Module{reqModule, globalModule},
	})
	providers := module.getRequest()
	fmt.Println(providers, module.DataProviders)
	require.Equal(t, 1, len(providers))
	require.Equal(t, Provide("req"), providers[0].Name)
}

func Test_appendProvider(t *testing.T) {
	module := NewModule(NewModuleOptions{
		Scope: Global,
	})

	provider := module.NewProvider(ProviderOptions{
		Name:  "test",
		Value: "test",
	})

	module.appendProvider(provider)
	require.Equal(t, 1, len(module.DataProviders))
}
