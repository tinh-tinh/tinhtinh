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

func ChildProvider(module core.Module) core.Provider {
	provider := module.NewProvider(core.ProviderOptions{
		Name:  "child",
		Value: "child",
	})
	return provider
}

func ChildModule(module core.Module) core.Module {
	childModule := module.New(core.NewModuleOptions{
		Scope:     core.Global,
		Providers: []core.Providers{ChildProvider},
		Exports:   []core.Providers{ChildProvider},
	})

	return childModule
}

func AppController(module core.Module) core.Controller {
	ctrl := module.NewController("test")
	ctrl.Get("", func(ctx core.Ctx) error {
		name := ctrl.Ref("child")
		return ctx.JSON(core.Map{
			"data": name,
		})
	})
	return ctrl
}

func AppModule() core.Module {
	module := core.NewModule(core.NewModuleOptions{
		Scope:       core.Global,
		Imports:     []core.Modules{ChildModule},
		Controllers: []core.Controllers{AppController},
	})

	return module
}

func Test_NewProvider(t *testing.T) {
	app := core.CreateFactory(AppModule)
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
	require.Equal(t, "child", res.Data)
}

func Test_FactoryProvider(t *testing.T) {
	rootModule := func(module core.Module) core.Module {
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

	childModule := func(module core.Module) core.Module {
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
		Imports: []core.Modules{rootModule, childModule},
	})

	require.Equal(t, "rootChild", module.Ref("child"))
}

func tenantModule() core.Module {
	const (
		TENANT  core.Provide = "TENANT"
		SERVICE core.Provide = "SERVICE"
	)
	type RequestProvider struct {
		Name string
	}
	service := func(module core.Module) core.Provider {
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

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			service := ctx.Ref(SERVICE).(*RequestProvider)
			return ctx.JSON(core.Map{
				"data": service.Name,
			})
		})

		return ctrl
	}

	tenantModule := func(module core.Module) core.Module {
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

	module := core.NewModule(core.NewModuleOptions{
		Imports:     []core.Modules{tenantModule},
		Controllers: []core.Controllers{controller},
		Providers:   []core.Providers{service},
	})

	return module
}

func Test_RequestProvider(t *testing.T) {
	app := core.CreateFactory(tenantModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
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

func Test_TransientProvider(t *testing.T) {
	type StructName struct {
		Name string
	}
	parentModule := func(module core.Module) core.Module {
		parent := module.New(core.NewModuleOptions{
			Scope: core.Global,
		})
		parent.NewProvider(core.ProviderOptions{
			Name: "parent",
			Value: &StructName{
				Name: "parent",
			},
		})
		parent.Export("parent")

		return parent
	}

	childModule := func(module core.Module) core.Module {
		child := module.New(core.NewModuleOptions{
			Scope: core.Transient,
		})
		child.NewProvider(core.ProviderOptions{
			Name: "child",
			Factory: func(param ...interface{}) interface{} {
				parent := param[0].(*StructName)
				return &StructName{
					Name: fmt.Sprintf("%vChild", parent.Name),
				}
			},
			Inject: []core.Provide{core.Provide("parent")},
		})
		child.Export("child")

		return child
	}

	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{parentModule, childModule},
	})

	// Singleton
	oldP, ok := module.Ref("parent").(*StructName)
	require.True(t, ok)
	require.Equal(t, "parent", oldP.Name)

	newP, ok := module.Ref("parent").(*StructName)
	require.True(t, ok)
	require.Equal(t, "parent", newP.Name)

	require.Same(t, oldP, newP)

	// Transient
	old, ok := module.Ref("child").(*StructName)
	require.True(t, ok)
	require.Equal(t, "parentChild", old.Name)

	newStr, ok := module.Ref("child").(*StructName)
	require.True(t, ok)
	require.Equal(t, "parentChild", newStr.Name)

	require.NotSame(t, old, newStr)
}

func Test_AutoNameProvider(t *testing.T) {
	type StructName struct {
		Name string
	}
	service := func(module core.Module) core.Provider {
		return module.NewProvider(&StructName{Name: "module"})
	}

	module := core.NewModule(core.NewModuleOptions{
		Providers: []core.Providers{service},
		Exports:   []core.Providers{service},
	})

	structName := core.Inject[StructName](module)
	require.Equal(t, "module", structName.Name)

	type Any struct {
		Version int
	}
	nilValue := core.Inject[Any](module)
	require.Nil(t, nilValue)
}

func BenchmarkRequestModule(b *testing.B) {
	app := core.CreateFactory(tenantModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
			require.Nil(b, err)
			req.Header.Set("x-tenant", "test")
			resp, err := testClient.Do(req)
			require.Nil(b, err)
			require.Equal(b, http.StatusOK, resp.StatusCode)
		}
	})
}
