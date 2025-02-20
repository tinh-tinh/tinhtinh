package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/router"
)

func Test_Module(t *testing.T) {
	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("auth")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "auth",
			})
		})

		return ctrl
	}

	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{userController},
		})

		return mod
	}

	postController := func(module core.Module) core.Controller {
		ctrl := module.NewController("documents")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "documents",
			})
		})

		return ctrl
	}

	postModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{postController},
		})

		return mod
	}

	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
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

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"auth"}`, string(data))

	resp, err = testClient.Get(testServer.URL + "/api/setting/documents")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"documents"}`, string(data))
}
