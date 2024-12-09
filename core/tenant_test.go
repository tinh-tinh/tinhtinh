package core_test

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func BenchmarkTenant(b *testing.B) {
	const (
		TENANCY core.Provide = "TENANCY"
		MAPPER  core.Provide = "MAPPER"
	)

	type mapper map[string]interface{}
	var mutex = sync.RWMutex{}

	forRoot := func(module *core.DynamicModule) *core.DynamicModule {
		tenantModule := module.New(core.NewModuleOptions{})

		tenantModule.NewProvider(core.ProviderOptions{
			Name:  MAPPER,
			Value: make(mapper),
		})
		tenantModule.Export(MAPPER)

		tenantModule.NewProvider(core.ProviderOptions{
			Scope: core.Request,
			Name:  TENANCY,
			Factory: func(param ...interface{}) interface{} {
				req := param[0].(*http.Request)
				tenantID := req.Header.Get("x-tenant")
				if tenantID == "" {
					tenantID = "master"
				}
				mutex.RLock()
				mapper, ok := param[1].(mapper)
				mutex.RUnlock()
				if !ok {
					return nil
				}

				if mapper[tenantID] == nil {
					mutex.Lock()
					mapper[tenantID] = tenantID
					mutex.Unlock()
				}
				return mapper[tenantID]
			},
			Inject: []core.Provide{core.REQUEST, MAPPER},
		})
		tenantModule.Export(TENANCY)

		return tenantModule
	}

	userController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("user")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Ref(TENANCY),
			})
		})

		return ctrl
	}

	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		user := module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{userController},
		})

		return user
	}

	appModule := func() *core.DynamicModule {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{forRoot, userModule},
		})

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req, err := http.NewRequest("GET", testServer.URL+"/api/user", nil)
			require.Nil(b, err)
			abc := rand.Intn(100)
			id := strconv.Itoa(abc)
			req.Header.Set("x-tenant", "test"+id)
			resp, err := testClient.Do(req)
			require.Nil(b, err)
			require.Equal(b, http.StatusOK, resp.StatusCode)

			data, err := io.ReadAll(resp.Body)
			require.Nil(b, err)

			var res Response
			err = json.Unmarshal(data, &res)
			require.Nil(b, err)
			require.Equal(b, "test"+id, res.Data)
		}
	})
}
