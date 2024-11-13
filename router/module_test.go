package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/router"
)

func Test_Module(t *testing.T) {
	userController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("auth")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{
			Controllers: []core.Controller{userController},
		})

		return mod
	}

	postController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("documents")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	postModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{
			Controllers: []core.Controller{postController},
		})

		return mod
	}

	appModule := func() *core.DynamicModule {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{
				router.Register(
					router.Options{
						Path: "setting",
						Children: []*router.RouteChildren{
							{
								Module: userModule,
							},
							{
								Module: postModule,
							},
						},
					},
					router.Options{
						Path:   "newnews",
						Module: postModule,
					},
				),
			},
		})

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/setting/auth")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/setting/documents")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
